package entity

import "time"

const (
	BannerShowStatus    string = "选择商家展示"
	BannerDeletedStatus string = "移除该商家展示位"
)

type BannerBody struct {
	Id        int64     `json:"id"`
	AppId     int64     `json:"app_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
