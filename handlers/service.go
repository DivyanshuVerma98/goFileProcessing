package handlers

import (
	"fmt"
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
		policy_no_list := []string{}
		for policy_no := range batch_data.MotorPolicy {
			_, exists := batch_data.Error[policy_no]
			if !exists {
				policy_no_list = append(policy_no_list, policy_no)
			}
		}
		// calling
		database.GetPolicyNo(policy_no_list)
		// destination <- batch_data
	}
}
