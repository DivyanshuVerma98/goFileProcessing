package structs

type PolicyDetails struct {
	DataMap map[string]map[string]string
	RowMap  map[string][]string
}

type ErrorDetails struct {
	MessageMap map[string]string
}

type BatchData struct {
	PolicyDetails PolicyDetails
	ErrorDetails  ErrorDetails
}

func (bd *BatchData) Initialize() {
	bd.PolicyDetails = PolicyDetails{
		DataMap: map[string]map[string]string{},
		RowMap:  map[string][]string{},
	}
	bd.ErrorDetails = ErrorDetails{
		MessageMap: map[string]string{},
	}
}

func (bd *BatchData) Copy() *BatchData {
	copy_data := &BatchData{
		PolicyDetails: bd.PolicyDetails,
		ErrorDetails:  bd.ErrorDetails,
	}

	return copy_data
}


type Manager struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type UserData struct {
	Privileges     []string `json:"privileges"`
	SystemRoles    []string `json:"system_roles"`
	ID             int      `json:"id"`
	Username       string   `json:"username"`
	Department     string   `json:"department"`
	BusinessRole   string   `json:"business_role"`
	LastName       string   `json:"lastname"`
	FirstName      string   `json:"firstname"`
	Pospid         string   `json:"pospid"`
	BusinessName   string   `json:"business_name"`
	Phone          string   `json:"phone"`
	Manager        Manager  `json:"manager"`
	Certifications []string `json:"certifications"`
	Manages        []string `json:"manages"`
	Subordinates   []string `json:"subordinates"`
}

type GetUserDataAPIResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    UserData  `json:"data"`
}

// To send upload_file API response
type FileUploadResponse struct {
	SucessCount        int    `json:"success_count"`
	ErrorCount         int    `json:"error_count"`
	CompleteReportLink string `json:"complete_report_link"`
	ErrorReportLink    string `json:"error_report_link"`
}

// To capture Upload API response
type UploadAPIResponse struct {
	Data    Data   `json:"data"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Data struct {
	ReferenceID string `json:"referenceid"`
}
