package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/DivyanshuVerma98/goFileProcessing/database"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/DivyanshuVerma98/goFileProcessing/utils"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func MotorService(w http.ResponseWriter, r *http.Request) {
	timer := time.Now()
	log.Println("START TIME ", timer)
	file, handler, err := r.FormFile("data_file")
	userData := r.Context().Value(constants.UserDataKey).(*structs.UserData)
	fmt.Println("userData", userData)
	if err != nil {
		log.Println("Error retrieving file: ", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	log.Println("File Name - ", handler.Filename)

	// Selecting mapping according to bussiness role
	var csvToModelMap map[string]string
	if userData.BusinessRole == constants.MakerBusinessRole {
		csvToModelMap = constants.MotorMakerCSVToModelMap
	} else if userData.BusinessRole == constants.EnricherBusinessRole {
		csvToModelMap = constants.MotorEnricherCSVToModelMap
	} else {
		csvToModelMap = constants.MotorApproverCSVToModelMap
	}
	//-------------------------------------------

	// batch_size, _ := strconv.Atoi(os.Getenv("MOTOR_BATCH_SIZE"))
	batch_size := 1000 // Not to be set more than 1000. Else, will cause DB query issues.
	batch_generator_chan, err := batchGenerator(&file, batch_size,
		csvToModelMap, userData, timer)
	if err != nil {
		SendResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	valid_batch_chan := validateBatch(batch_generator_chan, timer)
	db_valid_batch_chan := dbValidationBatch(valid_batch_chan, userData, timer)

	response := structs.FileUploadResponse{}
	// --- Creating Reports ---------------------
	// id := uuid.New()
	completeReportFileName := "complete-report.csv" //fmt.Sprintf(constants.MotorCompleteReportFileName, id.String())
	errorReportFileName := "error-report.csv"       //fmt.Sprintf(constants.MotorErrorReportFileName, id.String())
	// defer os.Remove(completeReportFileName)
	// defer os.Remove(errorReportFileName)
	createReport(db_valid_batch_chan, &completeReportFileName,
		&errorReportFileName, &response, timer)
	log.Println("Done -->", time.Since(timer))
	// ------------------------------------------
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		// Uploading file on S3
		time.Sleep(time.Duration(5))
		response.CompleteReportLink = "complete-report.link"
	}()
	go func() {
		defer wg.Done()
		// Uploading file on S3
		time.Sleep(time.Duration(5))
		response.ErrorReportLink = "error-report.link"
	}()
	wg.Wait()
	// --- Creating entry in FileInfo model -----
	fileStatus := constants.Pass
	if response.SucessCount == 0 && response.ErrorCount > 0 {
		fileStatus = constants.Fail
	} else if response.SucessCount > 0 && response.ErrorCount > 0 {
		fileStatus = constants.PartialSuccess
	}
	fileInfomation := structs.FileInformation{
		ID:                    uuid.New(),
		Filename:              handler.Filename,
		Status:                fileStatus,
		ProductName:           constants.Motor,
		NumberOfPolicies:      (response.SucessCount + response.ErrorCount),
		NumberOfSuccess:       response.SucessCount,
		NumberOfFailure:       response.ErrorCount,
		TotalPremiumOfSuccess: response.TotalPremiumOfSuccess,
		TotalPremiumOfFailure: response.TotalPremiumOfFailure,
		CompleteReport:        response.CompleteReportLink,
		ErrorReport:           response.ErrorReportLink,
		BusinessRole:          userData.BusinessRole,
		CreatedBy:             userData.Username,
		UpdatedBy:             userData.Username,
	}
	result := database.DB.Create(&fileInfomation)
	log.Println("Result of creating entry in FileInformation model -> ", result)
	if result.Error != nil {
		log.Println("Error -->", result.Error.Error())
		SendResponse(w, result.Error.Error(), http.StatusInternalServerError, nil)
		return
	}
	//-------------------------------------------
	SendResponse(w, "Success", http.StatusOK, response)
}

func batchGenerator(file *multipart.File, batch_size int,
	csvToModelMap map[string]string, userData *structs.UserData,
	timer time.Time) (<-chan *structs.MotorBatchData, error) {
	log.Println("Inside BatchGenerator")
	log.Println("Starting batchGenerator", time.Since(timer))
	generator_chan := make(chan *structs.MotorBatchData)
	csv_reader := csv.NewReader(*file)
	// --- Reading and validating CSV headers ---
	headers, err := csv_reader.Read()
	if err != nil {
		log.Println("Error reading file: ", err)
		return nil, err

	}
	err = utils.ValidateHeaders(headers, csvToModelMap)
	if err != nil {
		log.Println("Invalid headers: ", err)
		return nil, err
	}
	//------------------------------------------
	go func() {
		defer close(generator_chan)
		// This will act as the key value for each policy
		// in BatchData.PolicyDetails.DataMap
		policyCount := 1
		// To keep track of batches
		batchCount := 0
		// To store the policy number from CSV,
		// making sure duplicate policy numbers doesn't exists
		policyNumberMap := map[string]bool{}
		batchData := structs.MotorBatchData{}
		batchData.Initialize()
		batchData.CsvHeaders = headers
		for {
			row, err := csv_reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Error reading CSV:", err)
				return
			}
			row_data := map[string]string{}
			for index, field_val := range row {
				key := csvToModelMap[headers[index]]
				field_val = strings.ReplaceAll(field_val, "%", "")
				row_data[key] = field_val
			}
			// Converting CSV row to struct
			var motorPolicy structs.MotorPolicy
			err = utils.MapToStruct(row_data, &motorPolicy)
			if err != nil {
				// catching type errors
				batchData.ErrorDetails.MessageMap[strconv.Itoa(policyCount)] = err.Error()
			}

			// Checking for redundant policy numbers
			_, exists := policyNumberMap[motorPolicy.PolicyNumber]
			if exists {
				batchData.ErrorDetails.MessageMap[strconv.Itoa(policyCount)] = "Redundant policy number found."
			} else {
				policyNumberMap[motorPolicy.PolicyNumber] = true
			}

			motorPolicy.TransactionType = strings.ToUpper(motorPolicy.TransactionType)
			motorPolicy.UpdatedBy = userData.Username
			motorPolicy.CreatedBy = userData.Username

			enricherStatus := strings.ToUpper(motorPolicy.EnricherStatus)
			enricherStatus = strings.TrimSpace(enricherStatus)
			if slices.Contains(constants.EnricherAllowedSuccessValues, enricherStatus) {
				motorPolicy.EnricherStatus = constants.SentToNextStage
				motorPolicy.ApproverStatus = constants.Pending
			} else if slices.Contains(constants.AllowedRejectedValues, enricherStatus) {
				motorPolicy.EnricherStatus = constants.Rejected
			} else {
				motorPolicy.EnricherStatus = constants.Pending
			}

			approverStatus := strings.ToUpper(motorPolicy.ApproverStatus)
			approverStatus = strings.TrimSpace(approverStatus)
			if slices.Contains(constants.AllowedRejectedValues, approverStatus) {
				motorPolicy.ApproverStatus = constants.Rejected
				motorPolicy.EnricherStatus = constants.Pushback
			} else if slices.Contains(constants.ApproverAllowedSuccesValues, approverStatus) {
				motorPolicy.EnricherStatus = constants.SentToNextStage
				motorPolicy.ApproverStatus = constants.Approved
			} else if approverStatus == constants.Pushback {
				motorPolicy.EnricherStatus = constants.Pushback
			} else {
				motorPolicy.ApproverStatus = constants.Pending
			}

			batchData.PolicyDetails.DataMap[strconv.Itoa(policyCount)] = motorPolicy
			batchData.PolicyDetails.RowMap[strconv.Itoa(policyCount)] = row
			policyCount += 1
			batchCount += 1
			if batchCount >= batch_size {
				batchCount = 0
				policyCount = 1
				copy := batchData.Copy()
				generator_chan <- copy
				batchData.Initialize()
			}
		}
		if batchCount > 0 {
			batchCount = 0
			copy := batchData.Copy()
			generator_chan <- copy
			batchData.Initialize()
		}
		log.Println("Ending batchGenerator", time.Since(timer))
		policyNumberMap = nil
	}()
	return generator_chan, nil
}

func validateBatch(sourceChan <-chan *structs.MotorBatchData, timer time.Time) <-chan *structs.MotorBatchData {
	log.Println("Inside validateBatch")
	log.Println("Starting validateBatch", time.Since(timer))
	generatorChan := make(chan *structs.MotorBatchData)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < runtime.NumCPU()-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for batchData := range sourceChan {
					validateFields(batchData)
					generatorChan <- batchData
				}
			}()
		}
		wg.Wait()
		close(generatorChan)
		log.Println("Ending validateBatch", time.Since(timer))
	}()
	return generatorChan
}

func validateFields(batchData *structs.MotorBatchData) {
	for key, motorPolicy := range batchData.PolicyDetails.DataMap {
		for _, fval := range constants.MotorMandatoryFields {
			val := reflect.ValueOf(motorPolicy).FieldByName(fval)
			if !val.IsValid() || val.Interface() == "" {
				batchData.ErrorDetails.MessageMap[key] = fval + " field can't be empty."
				continue
			}
		}
		err := utils.ValidateTransactionType(motorPolicy.TransactionType)
		if err != nil {
			batchData.ErrorDetails.MessageMap[key] = err.Error()
			continue
		}

		insurerName := motorPolicy.InsurerName
		insurerName = strings.TrimSpace(insurerName)
		insurerName = strings.ToLower(insurerName)

		product := motorPolicy.Product
		product = strings.TrimSpace(product)
		product = strings.ToUpper(product)
		insurerMandatoryFields := constants.MotorInsurerJsonData[insurerName][product]
		for _, fval := range insurerMandatoryFields {
			val := reflect.ValueOf(motorPolicy).FieldByName(fval)
			if !val.IsValid() || val.Interface() == "" {
				batchData.ErrorDetails.MessageMap[key] = fval + " can't be empty."
				continue
			}
		}

	}
}

func dbValidationBatch(sourceChan <-chan *structs.MotorBatchData,
	userData *structs.UserData, timer time.Time) <-chan *structs.MotorBatchData {
	log.Println("Inside dbValidationBatch")
	log.Println("Starting dbValidationBatch", time.Since(timer))
	generatorChan := make(chan *structs.MotorBatchData)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < runtime.NumCPU()-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for batchData := range sourceChan {
					dbValidation(batchData, userData, timer)
					generatorChan <- batchData
				}
			}()
		}
		wg.Wait()
		close(generatorChan)
		log.Println("Ending dbValidationBatch", time.Since(timer))
	}()
	return generatorChan
}

func dbValidation(batchData *structs.MotorBatchData,
	userData *structs.UserData, timer time.Time) {
	// fmt.Println("START TIME ", start)
	policyNoList := []string{}
	policyNoKeyMap := map[string]string{}
	bulkCreatePolicies := []structs.MotorPolicy{}
	bulkUpdatePolicies := []structs.MotorPolicy{}
	for key, motorPolicy := range batchData.PolicyDetails.DataMap {
		_, exists := batchData.ErrorDetails.MessageMap[key]
		if !exists {
			policyNoList = append(policyNoList, motorPolicy.PolicyNumber)
			policyNoKeyMap[motorPolicy.PolicyNumber] = key
		}
	}
	log.Println("Getting data ...", time.Since(timer))
	var policies []structs.MotorPolicy

	// Find the latest entry for each policy_number in the policyNoList
	result := database.DB.Select("DISTINCT ON (policy_number) *").
		Where("policy_number IN (?)", policyNoList).
		Order("policy_number, created_at DESC").
		Find(&policies)

	if result.Error != nil {
		log.Println("Get query ERROR -->", result.Error)
	}
	log.Println("Got data ...", time.Since(timer))

	for _, policy := range policies {
		key := policyNoKeyMap[policy.PolicyNumber]
		// Getting the Motor Policy using the polcyNo and batchData key map
		motorPolicy := batchData.PolicyDetails.DataMap[policyNoKeyMap[policy.PolicyNumber]]
		if userData.BusinessRole == constants.MakerBusinessRole {
			batchData.ErrorDetails.MessageMap[key] = "Policy already exists."
			delete(policyNoKeyMap, policy.PolicyNumber)
		} else {
			motorPolicy.UpdatedBy = userData.Username
			motorPolicy.CreatedBy = policy.CreatedBy
			bulkUpdatePolicies = append(bulkUpdatePolicies, motorPolicy)
			delete(policyNoKeyMap, policy.PolicyNumber)
		}
	}
	for _, kval := range policyNoKeyMap {
		// Only Maker is allowed to create new entries
		if userData.BusinessRole != constants.MakerBusinessRole {
			batchData.ErrorDetails.MessageMap[kval] = fmt.Sprintf("%s can only update entries from existing policies",
				userData.BusinessRole)
		} else {
			_, exists := batchData.ErrorDetails.MessageMap[kval]
			if exists {
				continue
			}
			motorPolicy := batchData.PolicyDetails.DataMap[kval]
			motorPolicy.ID = uuid.New()
			motorPolicy.TotalOutAmount = (motorPolicy.OutAmount + motorPolicy.TotalOutAmount)
			bulkCreatePolicies = append(bulkCreatePolicies, motorPolicy)
		}
	}
	// --- Creating new entries in DB -----------
	if len(bulkCreatePolicies) > 0 {
		result := database.DB.Create(&bulkCreatePolicies)
		log.Println("Bulk Create result", result)
		if result.Error != nil {
			log.Println("Bulk Create Error -->", result.Error.Error())
			for key := range batchData.PolicyDetails.DataMap {
				batchData.ErrorDetails.MessageMap[key] = result.Error.Error()
			}
		}
	}
	//-------------------------------------------
	// --- Updating existing entries in DB ------
	if len(bulkUpdatePolicies) > 0 {
		result := database.DB.Save(bulkUpdatePolicies)
		log.Println("Bulk Update result", result)
		if result.Error != nil {
			log.Println("Bulk Update Error -->", result.Error.Error())
			for key := range batchData.PolicyDetails.DataMap {
				batchData.ErrorDetails.MessageMap[key] = result.Error.Error()
			}
		}
	}
	//-------------------------------------------
	log.Println("Batch Done", time.Since(timer))
}

func createReport(sourceChan <-chan *structs.MotorBatchData, completeReportFileNamem *string,
	errorReportFileName *string, reportResponse *structs.FileUploadResponse, timer time.Time) error {
	log.Println("Starting createReport", time.Since(timer))
	complete_report, err := os.Create(*completeReportFileNamem)
	if err != nil {
		return err
	}
	complete_report_writer := csv.NewWriter(complete_report)
	defer complete_report_writer.Flush()
	error_report, err := os.Create(*errorReportFileName)
	if err != nil {
		return err
	}
	error_report_writer := csv.NewWriter(error_report)
	defer error_report_writer.Flush()
	batchCount := 0
	errorCount := 0
	successCount := 0
	totalPremiumOfSuccess := 0.0
	totalPremiumOfFailure := 0.0
	for batchData := range sourceChan {
		for key, row := range batchData.PolicyDetails.RowMap {
			// -- Adding headers in report csv --
			if batchCount == 0 {
				batchCount = 1
				headers := batchData.CsvHeaders
				headers = append(headers, "Status", "Remark")
				error_report_writer.Write(headers)
				complete_report_writer.Write(headers)
			}
			// ----------------------------------
			errorMsg, exists := batchData.ErrorDetails.MessageMap[key]
			policyDetails := batchData.PolicyDetails.DataMap[key]
			if exists {
				errorCount += 1
				totalPremiumOfFailure += policyDetails.Premium
				row = append(row, constants.Failure, errorMsg)
				error_report_writer.Write(row)
				complete_report_writer.Write(row)
			} else {
				successCount += 1
				totalPremiumOfSuccess += policyDetails.Premium
				row = append(row, constants.Success)
				complete_report_writer.Write(row)
			}
		}

	}
	reportResponse.SucessCount = successCount
	reportResponse.ErrorCount = errorCount
	reportResponse.TotalPremiumOfSuccess = totalPremiumOfSuccess
	reportResponse.TotalPremiumOfFailure = totalPremiumOfFailure
	log.Println("Ending createReport", time.Since(timer))
	return nil
}
