package entity

import "time"

type News struct {
	Id        int64     `json:"id"`
	Pictures  string    `json:"pictures"`
	Name      string    `json:"name"`
	Describe  string    `json:"describe"`
	Status    Status    `json:"status"` // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsResp struct {
	Id        int64     `json:"id"`
	Pictures  []string  `json:"pictures"`
	PicUrl    []string  `json:"pic_url"`
	Name      string    `json:"name"`
	Describe  string    `json:"describe"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewsListBody struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
	Status   Status `json:"status"` // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	Option   int64  `json:"option"`
}

type NewsListResp struct {
	Total int64          `json:"total"`
	Apps  []NewsListBody `json:"list"`
}
