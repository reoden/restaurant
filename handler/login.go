package handler

import (
	"log"
	"restaurant/application"
	"restaurant/domain"
	"restaurant/entity"

	"github.com/gin-gonic/gin"
)

// LoginHandler ...
type LoginHandler struct {
	service application.LoginService
	account application.AccountService
}

// NewLoginHandler ...
func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		service: *application.NewLoginService(),
		account: *application.NewAccountService(),
	}
}

func (l *LoginHandler) SendCode(ctx *gin.Context) (interface{}, error) {
	var param domain.Login

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return l.service.SendCode(ctx, param)
}

func (l *LoginHandler) UserLogin(ctx *gin.Context) (interface{}, error) {
	var param domain.Login

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	user, err := l.account.GetAccountByPhone(ctx, param.Phone)
	if err != nil {
		log.Printf("GetAccount_ByPhone_err: %+v", err)
		return nil, err
	}

	// 用户不存在就创建
	if user == nil {
		ac := entity.Account{
			Phone: param.Phone,
		}
		_, err = l.account.Create(ctx, &ac)
	}

	if err != nil {
		return nil, err
	}

	return l.service.UserLogin(ctx, param)
}

func (l *LoginHandler) AdminLogin(ctx *gin.Context) (interface{}, error) {
	var param domain.Login

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	user, err := l.account.GetAccountByPhone(ctx, param.Phone)
	if err != nil {
		log.Printf("GetAccount_ByPhone_err: %+v", err)
		return nil, err
	}
	if user == nil {
		ac := entity.Account{
			Phone: param.Phone,
		}
		_, err = l.account.Create(ctx, &ac)
	}

	if err != nil {
		return nil, err
	}

	return l.service.AdminLogin(ctx, param)
}
