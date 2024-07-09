package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsTotalHandler struct {
	service     application.NewsService
	listService application.NewsShowListService
}

// NewNewsHandler ...
func NewNewsTotalHandler() *NewsTotalHandler {
	return &NewsTotalHandler{
		service:     *application.NewNewsService(),
		listService: *application.NewNewsShowListService(),
	}
}

func (ch *NewsTotalHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.NewsResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	app, err := ch.service.Create(ctx, param, userId)
	if err != nil {
		return nil, err
	}

	data := entity.NewsShowListBody{
		AppId: app.Id,
	}

	_, err = ch.listService.Create(ctx, userId, &data)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (ch *NewsTotalHandler) Delete(ctx *gin.Context) (interface{}, error) {
	accountID, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	err = ch.service.Delete(ctx, accountID, id)
	if err != nil {
		return nil, err
	}

	err = ch.listService.Delete(ctx, accountID, id)
	return nil, err
}

func (ch *NewsTotalHandler) List(ctx *gin.Context) (interface{}, error) {
	appResps, err := ch.listService.List(ctx)
	if err != nil {
		return nil, err
	}

	apps := appResps.Apps
	appIds := make([]int64, 0)
	for _, app := range apps {
		appIds = append(appIds, app.AppId)
	}

	if len(appIds) == 0 {
		return nil, nil
	}

	params, err := ch.service.GetApps(ctx, appIds)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"apps":  params,
		"total": appResps.Total,
	}, nil
}
