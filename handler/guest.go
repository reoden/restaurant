package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GuestHandler struct {
	service application.GuestService
}

// NewGuestHandler ...
func NewGuestHandler() *GuestHandler {
	return &GuestHandler{
		service: *application.NewGuestService(),
	}
}

// Create ...
func (ah *GuestHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.GuestResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return ah.service.Create(ctx, param, userId)
}

func (ah *GuestHandler) List(ctx *gin.Context) (interface{}, error) {
	return ah.service.List(ctx)
}

// Get ...
func (ah *GuestHandler) Get(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.Get(ctx, userId, id)
}

// Get ...
func (ah *GuestHandler) Read(ctx *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.ReadApp(ctx, id)
}

// Delete ...
func (ah *GuestHandler) Delete(ctx *gin.Context) (interface{}, error) {
	accountID, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return nil, ah.service.Delete(ctx, accountID, id)
}
