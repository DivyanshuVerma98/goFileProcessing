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
)

func MotorService(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("data_file")
	userData := r.Context().Value(UserDataKey)
	fmt.Println("userData", userData)
	if err != nil {
		log.Println("Error retrieving file: ", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	log.Println("File Name - ", handler.Filename)
	defer file.Close()
	csv_reader := csv.NewReader(file)
	headers, err := csv_reader.Read()
	if err != nil {
		log.Println("Error reading file: ", err)
		SendResponse(w, "Error reading file", http.StatusBadRequest, nil)
		return

	}
	err = utils.ValidateHeaders(headers, constants.MotorMakerCSVToModelMap)
	if err != nil {
		log.Println("Invalid headers: ", err)
		SendResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	// Reset file pointer to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("Error resetting file pointer:", err)
		SendResponse(w, "Error processing file", http.StatusInternalServerError, nil)
		return
	}

	// batch_size, _ := strconv.Atoi(os.Getenv("MOTOR_BATCH_SIZE"))
	batch_size := 5000
	batch_generator_chan := batchGenerator(&file, batch_size,
		constants.MotorMakerCSVToModelMap)
	valid_batch_chan := validateBatch(batch_generator_chan)
	db_valid_batch_chan := dbValidationBatch(valid_batch_chan)

	response := structs.FileUploadResponse{}
	createReport(db_valid_batch_chan, &response)
	mediaAuthToken := os.Getenv("MEDIA_URL_AUTH_TOKEN")
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		response.CompleteReportLink, err = UploadFile("complete_report.csv", mediaAuthToken)
		if err != nil {
			SendResponse(w, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	}()
	go func() {
		defer wg.Done()
		response.ErrorReportLink, err = UploadFile("error_report.csv", mediaAuthToken)
		if err != nil {
			SendResponse(w, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	}()
	wg.Wait()
	SendResponse(w, "Success", http.StatusOK, response)
}

func batchGenerator(file *multipart.File, batch_size int, csv_to_model_map map[string]string) <-chan *structs.MotorBatchData {
	log.Println("Inside BatchGenerator")
	generator_chan := make(chan *structs.MotorBatchData)
	csv_reader := csv.NewReader(*file)
	headers, _ := csv_reader.Read()
	go func() {
		defer close(generator_chan)
		// This will act as the key value for each policy
		// in BatchData.PolicyDetails.DataMap
		policyCount := 1
		// To keep track of batches
		batchCount := 0
		batchData := structs.MotorBatchData{}
		batchData.Initialize()
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
			field_list := []string{}
			for index, field_val := range row {
				key := csv_to_model_map[headers[index]]
				field_val = strings.ReplaceAll(field_val, "%", "")
				row_data[key] = field_val
				field_list = append(field_list, field_val)
			}
			var motorPolicy structs.MotorPolicy
			err = utils.MapToStruct(row_data, &motorPolicy)
			if err != nil {
				batchData.ErrorDetails.MessageMap[strconv.Itoa(policyCount)] = err.Error()
			}
			motorPolicy.TransactionType = strings.ToUpper(motorPolicy.TransactionType)
			motorPolicy.ApproverStatus = constants.Pending
			motorPolicy.EnricherStatus = constants.Pending

			batchData.PolicyDetails.DataMap[strconv.Itoa(policyCount)] = motorPolicy
			batchData.PolicyDetails.RowMap[strconv.Itoa(policyCount)] = field_list
			policyCount += 1
			// batch_data.ValidList = append(batch_data.ValidList, row_data)
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
	}()
	return generator_chan
}

func validateBatch(sourceChan <-chan *structs.MotorBatchData) <-chan *structs.MotorBatchData {
	log.Println("Inside validateBatch")
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
		// err = utils.IsValidDateFormat(motor_policy.BookingDate.String())
		// if err != nil {
		// 	batchData.ErrorDetails.MessageMap[key] = "Date format error. Invalid format - " + motor_policy.BookingDate.String()
		// }

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

func dbValidationBatch(sourceChan <-chan *structs.MotorBatchData) <-chan *structs.MotorBatchData {
	log.Println("Inside dbValidationBatch")
	generatorChan := make(chan *structs.MotorBatchData)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < runtime.NumCPU()-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for batchData := range sourceChan {
					dbValidation(batchData)
					generatorChan <- batchData
				}
			}()
		}
		wg.Wait()
		close(generatorChan)
	}()
	return generatorChan
}

func dbValidation(batchData *structs.MotorBatchData) {
	start := time.Now()
	// fmt.Println("START TIME ", start)
	policyNoList := []string{}
	policyNoKeyMap := map[string]string{}
	for key, motorPolicy := range batchData.PolicyDetails.DataMap {
		policyNoList = append(policyNoList, motorPolicy.PolicyNumber)
		policyNoKeyMap[motorPolicy.PolicyNumber] = key
	}
	// For Maker Flow
	policies, _ := database.GetMotorData("policy_number", policyNoList)
	// fmt.Println("TIME TO GET DATA ", time.Since(start))
	for _, policy := range policies {
		key := policyNoKeyMap[policy.PolicyNumber]
		if policy.ApproverStatus == constants.Approved {
			batchData.ErrorDetails.MessageMap[key] = "Policy already " + constants.Approved
		} else {
			batchData.ErrorDetails.MessageMap[key] = "Policy already exists."
		}
	}
	// fmt.Println("AFTER LOGIC ", time.Since(start))
	bulkCreatePolicies := []structs.MotorPolicy{}
	for key, motorPolicy := range batchData.PolicyDetails.DataMap {
		_, exists := batchData.ErrorDetails.MessageMap[key]
		if !exists {
			bulkCreatePolicies = append(bulkCreatePolicies, motorPolicy)
		}
	}

	err := database.MotorBulkCreate(bulkCreatePolicies)
	if err != nil {
		for key := range batchData.PolicyDetails.DataMap {
			batchData.ErrorDetails.MessageMap[key] = err.Error()
		}
	}
	fmt.Println("END TIME ", time.Since(start))
}

func createReport(sourceChan <-chan *structs.MotorBatchData, reportResponse *structs.FileUploadResponse) error {
	complete_report, err := os.Create("complete_report.csv")
	if err != nil {
		return err
	}
	complete_report_writer := csv.NewWriter(complete_report)
	defer complete_report_writer.Flush()
	error_report, err := os.Create("error_report.csv")
	if err != nil {
		return err
	}
	error_report_writer := csv.NewWriter(error_report)
	defer error_report_writer.Flush()
	errorCount := 0
	successCount := 0
	for batchData := range sourceChan {
		for key, row := range batchData.PolicyDetails.RowMap {
			errorMsg, exists := batchData.ErrorDetails.MessageMap[key]
			if exists {
				errorCount += 1
				row = append(row, constants.Failure, errorMsg)
				error_report_writer.Write(row)
				complete_report_writer.Write(row)
			} else {
				successCount += 1
				row = append(row, constants.Success)
				complete_report_writer.Write(row)
			}
		}

	}
	reportResponse.SucessCount = successCount
	reportResponse.ErrorCount = errorCount
	return nil
}
