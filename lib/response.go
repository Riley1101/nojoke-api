package lib

type Pagination struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Pagination
}

func NewResponse(status int, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
