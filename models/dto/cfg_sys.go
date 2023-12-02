package dto

import "time"

type CfgSys struct {
	BankId   string    `json:"bankId"`
	Name     string    `json:"name"`
	LastDate time.Time `json:"lastDate"`
	OpenDate time.Time `json:"openDate"`
	NextDate time.Time `json:"nextDate"`
	IpSource string    `json:"ipSource"`
	FlgSso   int       `json:"flgSso"`
}
