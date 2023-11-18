package lib

func PaginateData[T any](data []T, limit int, page int, total int) []T {
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
