package lib

import (
	"strconv"

	"github.com/gookit/validate"
)

func PaginateData[T interface{}](data []T, limit int, page int, total int) []T {
	tmp := []T{}
	if page*limit > total {
		return tmp
	}
	if page > 1 {
		data = data[(page-1)*limit:]
	}
	for i, v := range data {
		if i == limit {
			break
		}
		tmp = append(tmp, v)
	}
	return tmp
}

func PaginationParams(limit string, page string, total string) (int, int, int, error) {
	if limit == "" {
		limit = "10"
	}
	if page == "" {
		page = "1"
	}
	if total == "" {
		total = "100"
	}

	limitInt, error := strconv.Atoi(limit)
	pageInt, error := strconv.Atoi(page)
	totalInt, error := strconv.Atoi(total)

	return limitInt, pageInt, totalInt, error
}

func ValidateForm[T interface{}](form T) (bool, string) {
	v := validate.Struct(form)
	if !v.Validate() {
		message := v.Errors.One()
		return false, message
	}
	return true, ""
}
