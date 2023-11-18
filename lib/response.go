package lib

type Pagination struct {
	total int
	limit int
	page  int
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(status int, message string, data interface{}, p Pagination) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
