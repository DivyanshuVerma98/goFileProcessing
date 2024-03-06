package structs

type MotorPolicy struct {
	TransactionType         string `json:"transaction_type"`
	RmCode                  string `json:"rm_code"`
	RmName                  string `json:"rm_name"`
	ChildId                 string `json:"child_id"`
	BookingDate             string `json:"booking_date"`
	InsurerName             string `json:"insurer_name"`
	InsuredName             string `json:"insured_name"`
	MajorCategorisation     string `json:"major_categorisation"`
	Product                 string `json:"product"`
	ProductType             string `json:"product_type"`
	PolicyNo                string `json:"policy_no"`
	PlanType                string `json:"plan_type"`
	Premium                 string `json:"premium"`
	NetPremium              string `json:"net_premium"`
	OD                      string `json:"od"`
	TP                      string `json:"tp"`
	CommissionablePremium   string `json:"commissionable_premium"`
	RegistrationNo          string `json:"registation_no"`
	RtoCode                 string `json:"rto_code"`
	State                   string `json:"state"`
	RtoCluster              string `json:"rto_cluster"`
	City                    string `json:"city"`
	InsurerBiff             string `json:"insurer_biff"`
	FuelType                string `json:"fuel_type"`
	CPA                     string `json:"cpa"`
	CC                      string `json:"cc"`
	GVW                     string `json:"gvw"`
	NcbType                 string `json:"nsb_type"`
	SeatingCapacity         string `json:"seating_capacity"`
	VehicleRegistrationYear string `json:"vehicle_registration_year"`
	DiscountPercentage      string `json:"discount_percentage"`
	Make                    string `json:"make"`
	Model                   string `json:"model"`
	UTR                     string `json:"utr"`
	UTRDate                 string `json:"utr_date"`
	UtrAmount               string `json:"utr_amount"`
	Slot                    string `json:"slot"`
	PaidOnIn                string `json:"paid_on_in"`
	TentativeInPercentage   string `json:"tentative_in_percentage"`
	TentativeInAmount       string `json:"tentative_in_amount"`
	PaidOnOut               string `json:"paid_on_out"`
	OutPercentage           string `json:"out_percentage"`
	OutAmount               string `json:"out_amount"`
	CoType                  string `json:"co_type"`
	Remarks                 string `json:"remarks"`
}

type BatchData struct{
	MotorPolicy map[string]MotorPolicy
	Error map[string]string

}

type Response struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}
