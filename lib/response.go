package lib

type Pagination struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type DataResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination `json:"pagination"`
}

func NewDataResponse(status int, message string, data interface{}) DataResponse {
	return DataResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Message: message,
	}
}
