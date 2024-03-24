package structs

import "time"

type MotorPolicy struct {
	TransactionType         string    `json:"transaction_type"`
	RmCode                  string    `json:"rm_code"`
	RmName                  string    `json:"rm_name"`
	ChildId                 string    `json:"child_id"`
	BookingDate             time.Time `json:"booking_date"`
	InsurerName             string    `json:"insurer_name"`
	InsuredName             string    `json:"insured_name"`
	MajorCategory           string    `json:"major_category"`
	Product                 string    `json:"product"`
	ProductType             string    `json:"product_type"`
	PolicyNo                string    `json:"policy_no"`
	PlanType                string    `json:"plan_type"`
	Premium                 float64   `json:"premium"`
	NetPremium              float64   `json:"net_premium"`
	OD                      float64   `json:"od"`
	TP                      float64   `json:"tp"`
	CommissionablePremium   float64   `json:"commissionable_premium"`
	RegistrationNo          string    `json:"registation_no"`
	RtoCode                 string    `json:"rto_code"`
	State                   string    `json:"state"`
	RtoCluster              string    `json:"rto_cluster"`
	City                    string    `json:"city"`
	InsurerBiff             string    `json:"insurer_biff"`
	FuelType                string    `json:"fuel_type"`
	CPA                     string    `json:"cpa"`
	CC                      string    `json:"cc"`
	GVW                     string    `json:"gvw"`
	NcbType                 string    `json:"nsb_type"`
	SeatingCapacity         int       `json:"seating_capacity"`
	VehicleRegistrationYear int       `json:"vehicle_registration_year"`
	DiscountPercentage      float64   `json:"discount_percentage"`
	Make                    string    `json:"make"`
	Model                   string    `json:"model"`
	CTG                     string    `json:"ctg"`
	IDV                     string    `json:"idv"`
	UniqueId                string    `json:"unique_id"`
	SumInsuredVal           string    `json:"sum_insured_val"`
	VehicleRegistrationDate time.Time `json:"vehicle_registration_date"`
	UTR                     string    `json:"utr"`
	UTRDate                 time.Time `json:"utr_date"`
	UTRAmount               int       `json:"utr_amount"`
	SlotPaymentBatch        int       `json:"slot_payment_batch"`
	PaidOnIn                string    `json:"paid_on_in"`
	TentativeInPercentage   float32   `json:"tentative_in_percentage"`
	TentativeInAmount       float64   `json:"tentative_in_amount"`
	PaidOnOut               string    `json:"paid_on_out"`
	OutPercentage           float32   `json:"out_percentage"`
	OutAmount               float64   `json:"out_amount"`
	TotalOutAmount          float64   `json:"total_out_amount"`
	CoType                  string    `json:"co_type"`
	Remarks                 string    `json:"remarks"`
	BUHead                  string    `json:"bu_head"`
	Manager                 string    `json:"manager"`
	EnricherStatus          string    `json:"enricher_status"`
	ApproverStatus          string    `json:"approver_status"`
	EnricherRemark          string    `json:"enricher_remark"`
	ApproverRemark          string    `json:"approver_remark"`
}

type MotorPolicyDetails struct {
	DataMap map[string]MotorPolicy
	RowMap  map[string][]string
}
type MotorBatchData struct {
	PolicyDetails MotorPolicyDetails
	ErrorDetails  ErrorDetails
}
