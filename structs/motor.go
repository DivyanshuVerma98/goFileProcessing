package structs

import (
	"encoding/json"
	"fmt"
	"time"
)

type MotorPolicy struct {
	ID                      string    `db:"id"`
	TransactionType         string    `db:"transaction_type"`
	RmCode                  string    `db:"rm_code"`
	RmName                  string    `db:"rm_name"`
	ChildID                 string    `db:"child_id"`
	BookingDate             time.Time `db:"booking_date"`
	InsurerName             string    `db:"insurer_name"`
	InsuredName             string    `db:"insured_name"`
	MajorCategory           string    `db:"major_category"`
	Product                 string    `db:"product"`
	ProductType             string    `db:"product_type"`
	PolicyNumber            string    `db:"policy_number"`
	PlanType                string    `db:"plan_type"`
	Premium                 float64   `db:"premium"`
	NetPremium              float64   `db:"net_premium"`
	OD                      float64   `db:"od"`
	TP                      float64   `db:"tp"`
	CommissionablePremium   float64   `db:"commissionable_premium"`
	RegistrationNo          string    `db:"registration_no"`
	RTOCode                 string    `db:"rto_code"`
	State                   string    `db:"state"`
	RTOCluster              string    `db:"rto_cluster"`
	City                    string    `db:"city"`
	InsurerBiff             string    `db:"insurer_biff"`
	FuelType                string    `db:"fuel_type"`
	CPA                     string    `db:"cpa"`
	CC                      string    `db:"cc"`
	GVW                     string    `db:"gvw"`
	NCBType                 string    `db:"ncb_type"`
	SeatingCapacity         int       `db:"seating_capacity"`
	VehicleRegistrationYear int       `db:"vehicle_registration_year"`
	DiscountInPercentage    float64   `db:"discount_in_percentage"`
	Make                    string    `db:"make"`
	Model                   string    `db:"model"`
	CTG                     string    `db:"ctg"`
	IDV                     string    `db:"idv"`
	UniqueId                string    `db:"unique_id"`
	SumInsuredVal           string    `db:"sum_insured_val"`
	VehicleRegistrationDate time.Time `db:"vehicle_registration_date"`
	UTR                     string    `db:"utr"`
	UTRDate                 time.Time `db:"utr_date"`
	UTRAmount               int       `db:"utr_amount"`
	SlotPaymentBatch        string    `db:"slot_payment_batch"`
	PaidOnIn                string    `db:"paid_on_in"`
	TentativeInPercentage   float32   `db:"tentative_in_percentage"`
	TentativeInAmount       float64   `db:"tentative_in_amount"`
	PaidOnOut               string    `db:"paid_on_out"`
	OutPercentage           float32   `db:"out_percentage"`
	OutAmount               float64   `db:"out_amount"`
	TotalOutAmount          float64   `db:"total_out_amount"`
	COType                  string    `db:"co_type"`
	Remarks                 string    `db:"remarks"`
	BUHead                  string    `db:"bu_head"`
	Manager                 string    `db:"manager"`
	EnricherStatus          string    `db:"enricher_status"`
	ApproverStatus          string    `db:"approver_status"`
	EnricherRemark          string    `db:"enricher_remark"`
	ApproverRemark          string    `db:"approver_remark"`
}

func (mp *MotorPolicy) Copy() *MotorPolicy {
	// Convert the MotorPolicy struct to JSON
	jsonData, err := json.Marshal(mp)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	// Create a new MotorPolicy struct and unmarshal the JSON data into it
	var newMP MotorPolicy
	err = json.Unmarshal(jsonData, &newMP)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return &newMP
}

type MotorPolicyDetails struct {
	DataMap map[string]MotorPolicy
	RowMap  map[string][]string
}
type MotorBatchData struct {
	PolicyDetails MotorPolicyDetails
	ErrorDetails  ErrorDetails
}

func (m *MotorBatchData) Initialize() {
	m.PolicyDetails = MotorPolicyDetails{
		DataMap: map[string]MotorPolicy{},
		RowMap:  map[string][]string{},
	}
	m.ErrorDetails = ErrorDetails{
		MessageMap: map[string]string{},
	}
}

func (m *MotorBatchData) Copy() *MotorBatchData {
	copy_data := &MotorBatchData{
		PolicyDetails: m.PolicyDetails,
		ErrorDetails:  m.ErrorDetails,
	}

	return copy_data
}
