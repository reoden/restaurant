package entity

import "time"

type GuestShowListBody struct {
	Id        int64     `json:"id"`
	AppId     int64     `json:"app_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type GuestShowListResp struct {
	Total int64               `json:"total"`
	Apps  []GuestShowListBody `json:"apps"`
}
