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
	"sync"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/DivyanshuVerma98/goFileProcessing/database"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/DivyanshuVerma98/goFileProcessing/utils"
)

func MotorService(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("data_file")
	w.Header().Set("Content-Type", "application/json")
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
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	is_valid, msg := utils.ValidateHeaders(headers,
		constants.MOTOR_MAKER_CSV_TO_MODEL_MAP)
	if !is_valid {
		log.Println("Invalid headers: ", msg)
		http.Error(w, "Invalid headers. "+msg, http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("Hello")

}

func MotorBatchGenerator(file *multipart.File) <-chan *structs.BatchData {
	generator_chan := make(chan *structs.BatchData)
	csv_reader := csv.NewReader(*file)
	headers, _ := csv_reader.Read()
	fmt.Println("HEADERS ", headers)
	for _, val := range headers {
		fmt.Println(val)
	}
	// go func() {
	// 	defer close(generator_chan)
	// 	for {
	// 		row, err := csv_reader.Read()
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			fmt.Println(row)
	// 		}
	// 	}

	// }()
	return generator_chan
}

func ValidateBatchData(source chan *structs.BatchData,
	destination chan *structs.BatchData, wait_group *sync.WaitGroup) {
	fmt.Println("Inside ValidateBatchData")
	defer wait_group.Done()
	defer close(destination)
	for batch_data := range source {
		for policy_no, row_data := range batch_data.MotorPolicy {
			if !utils.IsValidDateFormat(row_data.BookingDate) {
				batch_data.Error[policy_no] = "This is the error"
			}
			fmt.Println("Validate", policy_no, "No Issues")
		}
		destination <- batch_data
	}
}

func QueryBatchData(source chan *structs.BatchData, wait_group *sync.WaitGroup) {
	fmt.Println("Inside QueryBatchData")
	defer wait_group.Done()
	database.CreateTable()
	for batch_data := range source {
		policy_list := []structs.MotorPolicy{}
		for _, policy_stuct := range batch_data.MotorPolicy {
			policy_list = append(policy_list, policy_stuct)
		}
		// calling
		// database.GetPolicyNo(policy_no_list)
		database.BulkInsert(policy_list)
		// destination <- batch_data
	}
}

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
