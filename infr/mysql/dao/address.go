package dao

type Address struct {
	Id      int    `json:"id" gorm:"primaryKey"`
	AppId   string `json:"app_id"`
	Address string `json:"address"`
	Enable  bool   `json:"enable"`
}

func (Address) TableName() string {
	return "address"
}
