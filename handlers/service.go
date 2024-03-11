package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/DivyanshuVerma98/goFileProcessing/database"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
	"github.com/DivyanshuVerma98/goFileProcessing/utils"
)

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

func UploadFile(filen_path string) string {
	url := "https://foundation-dev.theblackswan.in/media"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, _ := os.Open(filen_path)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", file.Name())
	if errFile1 != nil {
		fmt.Println(errFile1)
		return ""
	}
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return ""
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Authorization", "Token 43b75686cd6166057d572b5019accf9561ddbf00")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var api_reponse structs.UploadAPIResponse
	err = json.Unmarshal(body, &api_reponse)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("api_reponse", api_reponse.Data.ReferenceID)
	return url + "/" + api_reponse.Data.ReferenceID

}
