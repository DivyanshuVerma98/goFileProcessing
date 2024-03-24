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

func (bd *BatchData) GetInterface() interface{} {
	return bd
}

func (bd *BatchData) Copy() *BatchData {
	copy_data := &BatchData{
		PolicyDetails: bd.PolicyDetails,
		ErrorDetails:  bd.ErrorDetails,
	}

	return copy_data
}

type Response struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

type UploadAPIResponse struct {
	Data    Data   `json:"data"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Data struct {
	ReferenceID string `json:"referenceid"`
}
