package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/DivyanshuVerma98/goFileProcessing/structs"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	BATCH_SIZE := 1
	COUNT := 0
	file, handler, err := r.FormFile("data_file")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	fmt.Println("file_name ->", handler.Filename)
	csv_reader := csv.NewReader(file)

	// Removing headers
	csv_reader.Read()

	validate_channel := make(chan *structs.BatchData, 10)
	db_operation_channel := make(chan *structs.BatchData, 10)
	// create_report_channel := make(chan *structs.BatchData, 10)
	var wait_group sync.WaitGroup
	// result := make(chan bool)
	for i := 0; i < runtime.NumCPU()-1; i++ {
		wait_group.Add(1)
		go ValidateBatchData(validate_channel, db_operation_channel, &wait_group)
	}

	wait_group.Add(1)
	go QueryBatchData(db_operation_channel, &wait_group)

	batch_data := structs.BatchData{
		MotorPolicy: map[string]structs.MotorPolicy{},
		Error:       make(map[string]string),
	}
	go func(){
		for {
			row, err := csv_reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Panic(err)
			}

			policy := structs.MotorPolicy{
				TransactionType:         row[0],
				RmCode:                  row[1],
				RmName:                  row[2],
				ChildId:                 row[3],
				BookingDate:             row[4],
				InsurerName:             row[5],
				InsuredName:             row[6],
				MajorCategorisation:     row[7],
				Product:                 row[8],
				ProductType:             row[9],
				PolicyNo:                row[10],
				PlanType:                row[11],
				Premium:                 row[12],
				NetPremium:              row[13],
				OD:                      row[14],
				TP:                      row[15],
				CommissionablePremium:   row[16],
				RegistrationNo:          row[17],
				RtoCode:                 row[18],
				State:                   row[19],
				RtoCluster:              row[20],
				City:                    row[21],
				InsurerBiff:             row[22],
				FuelType:                row[23],
				CPA:                     row[24],
				CC:                      row[25],
				GVW:                     row[26],
				NcbType:                 row[27],
				SeatingCapacity:         row[28],
				VehicleRegistrationYear: row[29],
				DiscountPercentage:      row[30],
				Make:                    row[31],
				Model:                   row[32],
				UTR:                     row[33],
				UTRDate:                 row[34],
				UtrAmount:               row[35],
				Slot:                    row[36],
				PaidOnIn:                row[37],
				TentativeInPercentage:   row[38],
				TentativeInAmount:       row[39],
				PaidOnOut:               row[40],
				OutPercentage:           row[41],
				OutAmount:               row[42],
				CoType:                  row[43],
				Remarks:                 row[44],
			}
			batch_data.MotorPolicy[policy.PolicyNo] = policy
			COUNT += 1
			if COUNT >= BATCH_SIZE {
				data_copy := structs.BatchData{
					MotorPolicy: batch_data.MotorPolicy,
					Error:       batch_data.Error,
				}
				validate_channel <- &data_copy
				batch_data = structs.BatchData{
					MotorPolicy: map[string]structs.MotorPolicy{},
					Error:       make(map[string]string),
				}
				COUNT = 0
			}
		}
		if COUNT != 0 {
			data_copy := structs.BatchData{
				MotorPolicy: batch_data.MotorPolicy,
				Error:       batch_data.Error,
			}
			validate_channel <- &data_copy

		}
		close(validate_channel)
	}()
	wait_group.Wait()
	// <-result
	// close(result)
	w.WriteHeader(http.StatusOK)
	response := structs.Response{
		Status:  http.StatusOK,
		Message: "Success",
	}
	json.NewEncoder(w).Encode(response)
}
