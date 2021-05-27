package model

import "time"

type Trace struct {
	Id         int64     `json:"id"`
	Uuid       string    `json:"uuid"`
	Type       int       `json:"type"`
	Ip         string    `json:"ip"`
	CreateTime time.Time `json:"create_time"`
}
