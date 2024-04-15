package structs

import (
	"time"

	"github.com/google/uuid"
)

// -- File Infomation Model ---------
type FileInformation struct {
	ID                    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Filename              string    `gorm:"size:225" json:"filename"`
	Status                string    `gorm:"size:255" json:"status"`
	ProductName           string    `gorm:"size:255" json:"product_name"`
	NumberOfPolicies      int       `gorm:"type:integer" json:"number_of_policies"`
	NumberOfSuccess       int       `gorm:"type:integer" json:"number_of_success"`
	NumberOfFailure       int       `gorm:"type:integer" json:"number_of_failure"`
	TotalPremiumOfSuccess float64   `gorm:"type:double precision" json:"total_premium_of_success"`
	TotalPremiumOfFailure float64   `gorm:"type:double precision" json:"total_premium_of_failure"`
	CompleteReport        string    `gorm:"size:225" json:"complete_report"`
	ErrorReport           string    `gorm:"size:225" json:"error_report"`
	BusinessRole          string    `gorm:"size:255" json:"business_role"`
	UpdatedBy             string
	CreatedBy             string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (FileInformation) TableName() string {
	return "fms_fileinformation"
}

// ----------------------------------

// - To capture Get User Details API response
type UserData struct {
	Privileges     []string `json:"privileges"`
	SystemRoles    []string `json:"system_roles"`
	ID             int      `json:"id"`
	Username       string   `json:"username"`
	Department     string   `json:"department"`
	BusinessRole   string   `json:"business_role"`
	LastName       string   `json:"lastname"`
	FirstName      string   `json:"firstname"`
}
// -------------------------------------------

// --- To send upload_file API response ------
type FileUploadResponse struct {
	SucessCount           int     `json:"success_count"`
	ErrorCount            int     `json:"error_count"`
	TotalPremiumOfSuccess float64 `json:"total_premium_of_success"`
	TotalPremiumOfFailure float64 `json:"total_premium_of_failure"`
	CompleteReportLink    string  `json:"complete_report_link"`
	ErrorReportLink       string  `json:"error_report_link"`
}

// -------------------------------------------
