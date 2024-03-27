package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/DivyanshuVerma98/goFileProcessing/structs"
)

func GetUserDetailsAPI(userToken string) (*structs.UserData, error) {
	log.Println("Calling GetUserDetails API")
	url := os.Getenv("GET_USER_DETAILS_URL")
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", userToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	var response structs.GetUserDataAPIResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("APiI response", response)
	return &response.Data, nil
}

func UploadFile(file_path string, userToken string) (string, error) {
	log.Println("Calling UploadFile API")
	url := os.Getenv("MEDIA_URL")
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, _ := os.Open(file_path)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("file", file.Name())
	if errFile1 != nil {
		log.Println("UploadFile API Error ->", errFile1)
		return "", errFile1
	}
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		log.Println("UploadFile API Error ->", errFile1)
		return "", errFile1
	}
	err := writer.Close()
	if err != nil {
		log.Println("UploadFile API Error ->", err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("UploadFile API Error | Creating Request ->", err)
		return "", err
	}
	req.Header.Add("Authorization", userToken)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println("UploadFile API Error | Calling API ->", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("UploadFile API Error | Reading body ->", err)
		return "", err
	}
	var api_reponse structs.UploadAPIResponse
	err = json.Unmarshal(body, &api_reponse)
	if err != nil {
		log.Println("UploadFile API Error | Unmarshal response body ->", err)
		return "", err
	}
	log.Println("api_reponse", api_reponse)
	return url + "/" + api_reponse.Data.ReferenceID, nil

}
