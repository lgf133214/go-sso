package model

import "time"

type APP struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	AppId      string    `json:"app_id"`
	Status     int       `json:"status"`
	Redirect   string    `json:"redirect"`
	CreateTime time.Time `json:"create_time"`
}
