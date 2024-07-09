package entity

import "time"

const (
	CodeSended  string = "已发送"
	CodeUsed    string = "已验证"
	CodeExpired string = "已失效"
)

type VerificationCode struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
