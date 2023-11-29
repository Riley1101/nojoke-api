package collections

type Collection struct {
	Id       int64  `json:"id"`
	CreateAt string `json:"create_at"`
	UserId   int64  `json:"user_id"`
}

type ProductCollection struct {
	Collection
	Products []uint8 `json:"products"`
}
