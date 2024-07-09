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

type GuestService struct {
	db      *domain.GuestRepo
	account AccountService
}

func NewGuestService() *GuestService {
	return &GuestService{
		db:      domain.NewGuestRepo(),
		account: *NewAccountService(),
	}
}

func (as *GuestService) GetApp(ctx context.Context, id int64) (*entity.Guest, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *GuestService) ReadApp(ctx context.Context, id int64) (*entity.GuestResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.guest2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *GuestService) GetApps(ctx context.Context, ids []int64) ([]*entity.GuestResp, error) {
	apps, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.GuestResp, 0)
	for _, app := range apps {
		param, err := as.guest2Resp(app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return params, nil
}

func (a *GuestService) guest2Resp(app *entity.Guest) (*entity.GuestResp, error) {
	param := &entity.GuestResp{
		Id:        app.Id,
		Name:      app.Name,
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

func (a *GuestService) appResp2Guest(app *entity.GuestResp) (*entity.Guest, error) {
	param := &entity.Guest{
		Id:        app.Id,
		Name:      app.Name,
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
func (a *GuestService) Create(ctx context.Context, param entity.GuestResp, userId int64) (*entity.GuestResp, error) {
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
		newApp, err := a.appResp2Guest(&param)
		if err != nil {
			return nil, err
		}

		newApp.Status = entity.StatusAccepted
		newApp.CreatedAt = app.CreatedAt

		tempApp, err := a.Update(ctx, newApp)
		if err != nil {
			return nil, err
		}

		return a.guest2Resp(tempApp)
	}

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	params := entity.Guest{
		Pictures: string(data),
		Name:     param.Name,
		Describe: param.Describe,
		Status:   entity.StatusAccepted,
	}

	params.CreatedAt = currentTime
	params.UpdatedAt = currentTime

	tempApp, err := a.db.Create(ctx, &params)
	if err != nil {
		return nil, err
	}

	ret, err := a.guest2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List 列表
func (a *GuestService) List(ctx context.Context) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.GuestResp, 0)
	for _, app := range apps {
		param, err := a.guest2Resp(&app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *GuestService) Update(ctx context.Context, param *entity.Guest) (*entity.Guest, error) {
	return as.db.Update(ctx, param)
}

func (a *GuestService) Delete(ctx context.Context, userId, id int64) error {
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

func (a *GuestService) Get(ctx context.Context, userId, id int64) (*entity.GuestResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	guest, err := a.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	param, err := a.guest2Resp(guest)
	if err != nil {
		return nil, err
	}

	return param, nil
}
