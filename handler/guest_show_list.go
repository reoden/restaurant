package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GuestShowListHandler struct {
	service application.GuestShowListService
}

// NewGuestShowListHandler ...
func NewGuestShowListHandler() *GuestShowListHandler {
	return &GuestShowListHandler{
		service: *application.NewGuestShowListService(),
	}
}

func (bh *GuestShowListHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.GuestShowListBody

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return bh.service.Create(ctx, userId, &param)
}

func (bh *GuestShowListHandler) Delete(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return nil, bh.service.Delete(ctx, userId, id)
}

func (bh *GuestShowListHandler) Update(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	type updateParam struct {
		AppIds []int64 `json:"app_ids"`
	}

	var params updateParam
	if err := ctx.ShouldBind(&params); err != nil {
		return "", err
	}

	return nil, bh.service.Update(ctx, userId, params.AppIds)
}

func (bh *GuestShowListHandler) List(ctx *gin.Context) (interface{}, error) {
	return bh.service.List(ctx)
}
