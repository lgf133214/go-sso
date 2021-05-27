package model

import (
	"time"
)

type User struct {
	Id int64 `json:"id" `

	Uuid     string `json:"uuid"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Salt string `json:"salt"`
	IP   string

	Status     int    `json:"status"`
	VerifyCode string `json:"verify_code"`

	CreateTime time.Time `json:"create_time"`
	ModifyTime time.Time `json:"modify_time"`
}
