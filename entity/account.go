package entity

import (
	"time"
)

type AccountStatus int

const (
	UserStatus AccountStatus = iota
	AdministratorStatus
)

type Account struct {
	Id        int64         `json:"id"`
	Phone     string        `json:"phone"`
	Status    AccountStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

type LoginInfo struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"username"`
}
