package structs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MotorPolicy struct {
	ID                      uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TransactionType         string       `json:"transaction_type"`
	RmCode                  string       `json:"rm_code"`
	RmName                  string       `json:"rm_name"`
	ChildID                 string       `json:"child_id"`
	BookingDate             sql.NullTime `gorm:"type:date" json:"booking_date"`
	InsurerName             string       `json:"insurer_name"`
	InsuredName             string       `json:"insured_name"`
	MajorCategory           string       `json:"major_category"`
	Product                 string       `json:"product"`
	ProductType             string       `json:"product_type"`
	PolicyNumber            string       `gorm:"index" json:"policy_number"`
	PlanType                string       `json:"plan_type"`
	Premium                 float64      `json:"premium"`
	NetPremium              float64      `json:"net_premium"`
	OD                      float64      `json:"od"`
	TP                      float64      `json:"tp"`
	CommissionablePremium   float64      `json:"commissionable_premium"`
	RegistrationNo          string       `json:"registration_no"`
	RTOCode                 string       `json:"rto_code"`
	State                   string       `json:"state"`
	RTOCluster              string       `json:"rto_cluster"`
	City                    string       `json:"city"`
	InsurerBiff             string       `json:"insurer_biff"`
	FuelType                string       `json:"fuel_type"`
	CPA                     string       `json:"cpa"`
	CC                      string       `json:"cc"`
	GVW                     string       `json:"gvw"`
	NCBType                 string       `json:"ncb_type"`
	SeatingCapacity         int          `json:"seating_capacity"`
	VehicleRegistrationYear int          `json:"vehicle_registration_year"`
	DiscountInPercentage    float64      `json:"discount_in_percentage"`
	Make                    string       `json:"make"`
	Model                   string       `json:"model"`
	CTG                     string       `json:"ctg"`
	IDV                     string       `gorm:"column:idv" json:"idv"`
	UniqueId                string       `json:"unique_id"`
	SumInsuredVal           string       `json:"sum_insured_val"`
	VehicleRegistrationDate sql.NullTime `gorm:"type:date" json:"vehicle_registration_date"`
	UTR                     string       `json:"utr"`
	UTRDate                 sql.NullTime `gorm:"type:date" json:"utr_date"`
	UTRAmount               int          `json:"utr_amount"`
	SlotPaymentBatch        string       `json:"slot_payment_batch"`
	PaidOnIn                string       `json:"paid_on_in"`
	TentativeInPercentage   float32      `json:"tentative_in_percentage"`
	TentativeInAmount       float64      `json:"tentative_in_amount"`
	PaidOnOut               string       `json:"paid_on_out"`
	OutPercentage           float32      `json:"out_percentage"`
	OutAmount               float64      `json:"out_amount"`
	TotalOutAmount          float64      `json:"total_out_amount"`
	COType                  string       `json:"co_type"`
	Remarks                 string       `json:"remarks"`
	BUHead                  string       `json:"bu_head"`
	Manager                 string       `json:"manager"`
	EnricherStatus          string       `json:"enricher_status"`
	ApproverStatus          string       `json:"approver_status"`
	EnricherRemark          string       `json:"enricher_remark"`
	ApproverRemark          string       `json:"approver_remark"`
	UpdatedBy               string
	CreatedBy               string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (MotorPolicy) TableName() string {
	return "dms_motorinsurance"
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

type ErrorDetails struct {
	MessageMap map[string]string
}

type MotorBatchData struct {
	PolicyDetails MotorPolicyDetails
	ErrorDetails  ErrorDetails
	CsvHeaders    []string
}

func (m *MotorBatchData) Initialize() {
	m.PolicyDetails = MotorPolicyDetails{
		DataMap: map[string]MotorPolicy{},
		RowMap:  map[string][]string{},
	}
	m.ErrorDetails = ErrorDetails{
		MessageMap: map[string]string{},
	}
	m.CsvHeaders = []string{}
}

func (m *MotorBatchData) Copy() *MotorBatchData {
	copy_data := &MotorBatchData{
		PolicyDetails: m.PolicyDetails,
		ErrorDetails:  m.ErrorDetails,
		CsvHeaders:    m.CsvHeaders,
	}

	return copy_data
}
