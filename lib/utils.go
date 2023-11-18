package lib

import (
	"strconv"
)

func PaginateData[T interface{}](data []T, limit int, page int, total int) []T {
	tmp := []T{}
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
