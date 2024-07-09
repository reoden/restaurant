package entity

import "time"

type Chef struct {
	Id        int64     `json:"id"`
	Pictures  string    `json:"pictures"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Describe  string    `json:"describe"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChefResp struct {
	Id        int64     `json:"id"`
	Pictures  []string  `json:"pictures"`
	PicUrl    []string  `json:"pic_url"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Describe  string    `json:"describe"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChefListBody struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Describe string `json:"describe"`
	Status   Status `json:"status"`
}

type ChefListResp struct {
	Total int64          `json:"total"`
	Apps  []ChefListBody `json:"list"`
}
