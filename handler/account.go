package handler

import (
	"errors"
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service application.AccountService
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{
		service: *application.NewAccountService(),
	}
}

func (ah *AccountHandler) Self(ctx *gin.Context) (interface{}, error) {
	accountId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	return ah.service.GetAccountSelf(ctx, accountId)
}

func (ah *AccountHandler) UpdatePhone(ctx *gin.Context) (interface{}, error) {
	accountId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	password, ok := ctx.GetPostForm("password")
	if !ok {
		return ah.service.UpdatePhone(ctx, accountId, password), nil
	}

	return nil, errors.New("修改电话号码失败")
}

func (ah *AccountHandler) Register(ctx *gin.Context) (interface{}, error) {
	var param entity.Account

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return ah.service.Create(ctx, &param)
}
