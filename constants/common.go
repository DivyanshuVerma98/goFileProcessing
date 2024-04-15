package constants

// To set value in Context
// Need to define a custom type for context key
type customType string

const UserDataKey customType = "UserData"

const (
	Motor = "motor"

	MakerBusinessRole    = "fdms_maker"
	EnricherBusinessRole = "fdms_enricher"
	ApproverBusinessRole = "fdms_approver"

	Primary     = "PRIMARY"
	Adjustment  = "ADJUSTMENT"
	Endorsement = "ENDORSEMENT"

	Pending         = "PENDING"
	SentToNextStage = "SENT TO NEXT STAGE"
	Rejected        = "REJECTED"
	Approved        = "APPROVED"
	Pushback        = "PUSHBACK"

	Pass           = "PASS"
	PartialSuccess = "PARTIAL SUCCESS"
	Fail           = "FAIL"

	TransactionType         = "TransactionType"
	RmCode                  = "RmCode"
	RmName                  = "RmName"
	ChildID                 = "ChildID"
	BookingDate             = "BookingDate"
	InsurerName             = "InsurerName"
	InsuredName             = "InsuredName"
	MajorCategory           = "MajorCategory"
	Product                 = "Product"
	ProductType             = "ProductType"
	PolicyNumber            = "PolicyNumber"
	PlanType                = "PlanType"
	Premium                 = "Premium"
	NetPremium              = "NetPremium"
	OD                      = "OD"
	TP                      = "TP"
	CommissionablePremium   = "CommissionablePremium"
	RegistrationNo          = "RegistrationNo"
	RTOCode                 = "RTOCode"
	State                   = "State"
	RTOCluster              = "RTOCluster"
	City                    = "City"
	InsurerBiff             = "InsurerBiff"
	FuelType                = "FuelType"
	CPA                     = "CPA"
	CC                      = "CC"
	GVW                     = "GVW"
	NCBType                 = "NCBType"
	SeatingCapacity         = "SeatingCapacity"
	VehicleRegistrationYear = "VehicleRegistrationYear"
	DiscountInPercentage    = "DiscountInPercentage"
	Make                    = "Make"
	Model                   = "Model"
	CTG                     = "CTG"
	IDV                     = "IDV"
	SumInsuredVal           = "SumInsuredVal"
	UniqueID                = "UniqueId"
	VehicleRegistrationDate = "VehicleRegistrationDate"
	UTR                     = "UTR"
	UTRDate                 = "UTRDate"
	UTRAmount               = "UTRAmount"
	SlotPaymentBatch        = "SlotPaymentBatch"
	PaidOnIn                = "PaidOnIn"
	TentativeInPercentage   = "TentativeInPercentage"
	TentativeInAmount       = "TentativeInAmount"
	PaidOnOut               = "PaidOnOut"
	OutPercentage           = "OutPercentage"
	OutAmount               = "OutAmount"
	COType                  = "COType"
	Remarks                 = "Remarks"
	BUHead                  = "BUHead"
	Manager                 = "Manager"

	EnricherStatus = "EnricherStatus"
	ApproverStatus = "ApproverStatus"
	EnricherRemark = "EnricherRemark"
	ApproverRemark = "ApproverRemark"

	Failure = "Failure"
	Success = "Success"
)

var EnricherAllowedSuccessValues = []string{SentToNextStage, "SEND TO NEXT STAGE", "NEXT STAGE"}
var ApproverAllowedSuccesValues = []string{Approved, "APPROVE", "SUCCESS"}
var AllowedRejectedValues = []string{Rejected, "REJECT"}
