package event

type OrderNotify struct {
	Id     int64 `json:"id"`
	Status int   `json:"status"`
}
