package entities

import "time"

type ApiCall struct {
	ID             int `gorm:"primary_key"`
	CreatedAt      time.Time
	Track          string `gorm:"type:varchar(50)"`
	Service        string `gorm:"type:varchar(20)"`
	Webtype        string `gorm:"type:varchar(50)"`
	Merchant       string `gorm:"type:varchar(30)"`
	Msisdn         string `gorm:"type:varchar(50);index"`
	URL            string
	Method         string `gorm:"type:varchar(10)"`
	RequestQuery   string
	RequestHeader  string
	RequestBody    string
	ResponseHeader string
	ResponseBody   string
	StatusCode     int
	Latency        string
	Error          string

	// telco purpose
	TransactionID string `gorm:"type:varchar(150);index"`
}
