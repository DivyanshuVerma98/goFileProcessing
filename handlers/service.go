package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/DivyanshuVerma98/goFileProcessing/utils"
)

func MotorService(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("data_file")
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
		SendResponse(w, "Error reading file", http.StatusBadRequest, map[string]string{})
		return

	}
	err = utils.ValidateHeaders(headers, constants.MotorMakerCSVToModelMap)
	if err != nil {
		log.Println("Invalid headers: ", err)
		SendResponse(w, err.Error(), http.StatusBadRequest, map[string]string{})
		return
	}
	// Reset file pointer to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Println("Error resetting file pointer:", err)
		SendResponse(w, "Error processing file", http.StatusInternalServerError, map[string]string{})
		return
	}
	response := []interface{}{}
	// batch_size, _ := strconv.Atoi(os.Getenv("MOTOR_BATCH_SIZE"))
	batch_size := 1
	batch_generator_chan := BatchGenerator(&file, batch_size,
		constants.MotorMakerCSVToModelMap)
	valid_batch_chan := ValidateBatchGenerator(batch_generator_chan,
		reflect.TypeOf(structs.MotorPolicy{}))
	for val := range valid_batch_chan {
		utils.ValidateMotorBatchData(val)
		response = append(response, *val)
	}
	SendResponse(w, "Success", http.StatusOK, response)
}

func ValidateBatchGenerator(sourceChan <-chan *structs.BatchData,
	resType reflect.Type) <-chan *interface{} {
	log.Println("Inside ValidateBatchGenerator")
	generatorChan := make(chan *interface{})
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < runtime.NumCPU()-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if resType == reflect.TypeOf(structs.MotorPolicy{}) {
					fmt.Println("GOT IT")
				}
				for batchData := range sourceChan {
					for key, policy_data := range batchData.PolicyDetails.DataMap {
						if !utils.IsValidDateFormat(policy_data["booking_data"]) {
							batchData.ErrorDetails.MessageMap[key] = "Date format error. Invalid format - " + policy_data["booking_data"]
						}
					}
					generatorChan <- batchData.GetInterface()
				}
			}()
		}
		wg.Wait()
		close(generatorChan)
	}()
	return generatorChan
}

func BatchGenerator(file *multipart.File, batch_size int, csv_to_model_map map[string]string) <-chan *structs.BatchData {
	log.Println("Inside BatchGenerator")
	generator_chan := make(chan *structs.BatchData)
	csv_reader := csv.NewReader(*file)
	headers, _ := csv_reader.Read()
	go func() {
		defer close(generator_chan)
		// This will act as the key value for each policy
		// in BatchData.PolicyDetails.DataMap
		policy_count := 1
		// To keep track of batches
		batch_count := 0
		batch_data := structs.BatchData{}
		batch_data.Initialize()
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
				row_data[key] = field_val
				field_list = append(field_list, field_val)
			}
			batch_data.PolicyDetails.DataMap[strconv.Itoa(policy_count)] = row_data
			batch_data.PolicyDetails.RowMap[strconv.Itoa(policy_count)] = field_list
			policy_count += 1
			// batch_data.ValidList = append(batch_data.ValidList, row_data)
			batch_count += 1
			if batch_count >= batch_size {
				batch_count = 0
				copy := batch_data.Copy()
				generator_chan <- copy
				batch_data.Initialize()
			}
		}
		if batch_count > 0 {
			batch_count = 0
			copy := batch_data.Copy()
			generator_chan <- copy
			batch_data.Initialize()
		}
	}()
	return generator_chan
}

// func ValidateBatchData(source chan *structs.BatchData,
// 	destination chan *structs.BatchData, wait_group *sync.WaitGroup) {
// 	fmt.Println("Inside ValidateBatchData")
// 	defer wait_group.Done()
// 	defer close(destination)
// 	for batch_data := range source {
// 		for policy_no, row_data := range batch_data.MotorPolicy {
// 			if !utils.IsValidDateFormat(row_data.BookingDate) {
// 				batch_data.Error[policy_no] = "This is the error"
// 			}
// 			fmt.Println("Validate", policy_no, "No Issues")
// 		}
// 		destination <- batch_data
// 	}
// }

// func QueryBatchData(source chan *structs.BatchData, wait_group *sync.WaitGroup) {
// 	fmt.Println("Inside QueryBatchData")
// 	defer wait_group.Done()
// 	database.CreateTable()
// 	for batch_data := range source {
// 		policy_list := []structs.MotorPolicy{}
// 		for _, policy_stuct := range batch_data.MotorPolicy {
// 			policy_list = append(policy_list, policy_stuct)
// 		}
// 		// calling
// 		// database.GetPolicyNo(policy_no_list)
// 		database.BulkInsert(policy_list)
// 		// destination <- batch_data
// 	}
// }

func GetUserDetailsAPI() string {
	fmt.Println("Calling GetUserDetails API")
	url := os.Getenv("GET_USER_DETAILS_URL")
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.Header.Add("Authorization", "Token d08066ddabd18bc701a9c2b514577aa247b33c25")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	type Result struct {
		Status  int
		Message string
		Data    map[string]interface{}
	}

	var result Result
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("Result", result)
	fmt.Println(result.Data["business_role"])
	fmt.Println(result.Data["firstname"])
	fmt.Println(result.Data["lastname"])
	return ""
}

func UploadFile(file_path string, result chan string) {
	fmt.Println("Calling UploadFile API")
	url := os.Getenv("MEDIA_URL")
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, _ := os.Open(file_path)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", file.Name())
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Token 43b75686cd6166057d572b5019accf9561ddbf00")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var api_reponse structs.UploadAPIResponse
	err = json.Unmarshal(body, &api_reponse)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("api_reponse", file_path, api_reponse)
	result <- url + "/" + api_reponse.Data.ReferenceID

}

// "BatchData":{
// 	"PolicyDetails": {
// 		"data_map": {
// 			"1": {
// 				"booking_Data": ""
// 			}
// 		}
// 	},
// 	"ErrorDetials": {
// 		""
// 	}
// }

// {
// 	"PolicyData":{
// 		"1":{
// 			"booking_date": "20/02/2023"
// 		}
// 	},
// 	"ErrorData":{
// 		"1": {
// 			"message": "Booking date is not valid"
// 		}
// 	}
// }
