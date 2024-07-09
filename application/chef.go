package application

import (
	"context"
	"encoding/json"
	"log"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type ChefService struct {
	db      *domain.ChefRepo
	account AccountService
}

func NewChefService() *ChefService {
	return &ChefService{
		db:      domain.NewChefRepo(),
		account: *NewAccountService(),
	}
}

func (as *ChefService) GetApp(ctx context.Context, id int64) (*entity.Chef, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *ChefService) ReadApp(ctx context.Context, id int64) (*entity.ChefResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.chef2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *ChefService) GetApps(ctx context.Context, ids []int64) ([]*entity.ChefResp, error) {
	apps, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.ChefResp, 0)
	for _, app := range apps {
		param, err := as.chef2Resp(app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return params, nil
}

func (a *ChefService) chef2Resp(app *entity.Chef) (*entity.ChefResp, error) {
	param := &entity.ChefResp{
		Id:        app.Id,
		Name:      app.Name,
		Address:   app.Address,
		Describe:  app.Describe,
		Status:    app.Status,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}

	bodyPics := make([]string, 0)
	var appPics []string
	err := json.Unmarshal([]byte(app.Pictures), &appPics)
	if err != nil {
		return nil, err
	}

	for _, pic := range appPics {
		signedUrl, _ := pkgs.SignedUrl(pic, false)
		bodyPics = append(bodyPics, signedUrl)
	}

	log.Println(bodyPics)
	param.PicUrl = bodyPics
	param.Pictures = appPics

	return param, nil
}

func (a *ChefService) appResp2Chef(app *entity.ChefResp) (*entity.Chef, error) {
	param := &entity.Chef{
		Id:        app.Id,
		Name:      app.Name,
		Address:   app.Address,
		Describe:  app.Describe,
		Status:    app.Status,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}

	pics, err := json.Marshal(app.Pictures)
	if err != nil {
		return nil, err
	}

	param.Pictures = string(pics)

	return param, nil
}

// Create 创建
func (a *ChefService) Create(ctx context.Context, param entity.ChefResp, userId int64) (*entity.ChefResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	currentTime := common.NowInLocal()

	app, err := a.db.Get(ctx, param.Id)
	if err != nil {
		return nil, err
	}

	if app != nil {
		newApp, err := a.appResp2Chef(&param)
		if err != nil {
			return nil, err
		}

		newApp.Status = entity.StatusAccepted
		newApp.CreatedAt = app.CreatedAt

		tempApp, err := a.Update(ctx, newApp)
		if err != nil {
			return nil, err
		}

		return a.chef2Resp(tempApp)
	}

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	params := entity.Chef{
		Pictures: string(data),
		Name:     param.Name,
		Address:  param.Address,
		Describe: param.Describe,
		Status:   entity.StatusAccepted,
	}

	params.CreatedAt = currentTime
	params.UpdatedAt = currentTime

	tempApp, err := a.db.Create(ctx, &params)
	if err != nil {
		return nil, err
	}

	ret, err := a.chef2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List 列表
func (a *ChefService) List(ctx context.Context) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.ChefResp, 0)
	for _, app := range apps {
		param, err := a.chef2Resp(&app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
	// appListBody := make([]entity.ChefListBody, total)
	// for id, val := range apps {
	// 	appListBody[id].Address = val.Address
	// 	appListBody[id].Describe = val.Describe
	// 	appListBody[id].Name = val.Name
	// 	appListBody[id].Status = val.Status
	// 	appListBody[id].Id = val.Id
	// }
	// return &entity.ChefListResp{
	// 	Total: total,
	// 	Apps:  appListBody,
	// }, nil
}

func (as *ChefService) Update(ctx context.Context, param *entity.Chef) (*entity.Chef, error) {
	return as.db.Update(ctx, param)
}

func (a *ChefService) Delete(ctx context.Context, userId, id int64) error {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	_, err = a.db.Get(ctx, id)
	if err != nil {
		return err
	}

	return a.db.Delete(ctx, id)
}

func (a *ChefService) Get(ctx context.Context, userId, id int64) (*entity.ChefResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	chef, err := a.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	param, err := a.chef2Resp(chef)
	if err != nil {
		return nil, err
	}

	return param, nil
}
