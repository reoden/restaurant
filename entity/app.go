package entity

import "time"

type Status int

const (
	StatusPending Status = iota
	StatusRefused
	StatusAccepted
	StatusDeleted
	StatusSaved
)

type Application struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Pictures    string    `json:"pictures"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Describe    string    `json:"describe"`
	Phone       string    `json:"phone"`
	PostCode    string    `json:"post_code"`
	PostName    string    `json:"post_name"`
	WorkBeginAt string    `json:"work_begin_at"`
	WorkEndAt   string    `json:"work_end_at"`
	HaveVege    int       `json:"have_vege"` //0 没有素食 1 有素食
	Status      Status    `json:"status"`    // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppResp struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Pictures    []string  `json:"pictures"`
	PicsUrl     []string  `json:"pics_url"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Describe    string    `json:"describe"`
	Phone       string    `json:"phone"`
	PostCode    string    `json:"post_code"`
	PostName    string    `json:"post_name"`
	WorkBeginAt string    `json:"work_begin_at"`
	WorkEndAt   string    `json:"work_end_at"`
	HaveVege    int       `json:"have_vege"` //0 没有素食 1 有素食
	Status      Status    `json:"status"`    // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ApplicationListBody struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Describe string `json:"describe"`
	Status   Status `json:"status"` // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	Option   int64  `json:"option"`
}

type ApplicationListResp struct {
	Total int64                 `json:"total"`
	Apps  []ApplicationListBody `json:"list"`
}
