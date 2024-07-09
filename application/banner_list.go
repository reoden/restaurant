package application

import (
	"context"
	"errors"
	"log"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type BannerService struct {
	db      *domain.BannerRepo
	as      AppService
	account AccountService
}

func NewBannerService() *BannerService {
	return &BannerService{
		db:      domain.NewBannerRepo(),
		as:      *NewAppService(),
		account: *NewAccountService(),
	}
}

func (br *BannerService) Create(ctx context.Context, userId int64, param *entity.BannerBody) (*entity.BannerBody, error) {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	app, err := br.as.GetApp(ctx, param.AppId)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, errors.New("当前商家尚未通过审核")
	}

	banner, err := br.db.GetByAppId(ctx, param.AppId)
	if err != nil {
		return nil, err
	}

	if banner != nil {
		return nil, errors.New("重复添加")
	}

	currentTime := common.NowInLocal()
	param.CreatedAt = currentTime
	param.Status = entity.BannerShowStatus
	return br.db.Create(ctx, param)
}

func (br *BannerService) CreateBanners(ctx context.Context, param []*entity.BannerBody) ([]*entity.BannerBody, error) {
	currentTime := common.NowInLocal()
	for i := range param {
		param[i].CreatedAt = currentTime
		param[i].Status = entity.BannerShowStatus
	}
	return br.db.CreateBanners(ctx, param)
}

func (br *BannerService) Delete(ctx context.Context, userId, id int64) error {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	accountStatus := user.Status
	if accountStatus != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	return br.db.Delete(ctx, id)
}

func (br *BannerService) List(ctx context.Context) (map[string]interface{}, error) {
	apps, total, err := br.db.List(ctx)
	if err != nil {
		return nil, err
	}

	appIds := make([]int64, 0)
	for _, app := range apps {
		appIds = append(appIds, app.AppId)
	}

	log.Println(appIds)

	if len(appIds) == 0 {
		return nil, nil
	}

	params, err := br.as.GetAppsOrder(ctx, appIds)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"apps":  params,
		"total": total,
	}, nil
}

func (br *BannerService) Update(ctx context.Context, userId int64, appIds []int64) error {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	banners, err := br.db.GetAll(ctx)
	if err != nil {
		return err
	}

	bannerIds := make([]int64, 0, len(banners))
	for _, banner := range banners {
		bannerIds = append(bannerIds, banner.Id)
	}

	bannerBodys := make([]*entity.BannerBody, 0)
	for _, appId := range appIds {
		bannerBodys = append(bannerBodys, &entity.BannerBody{
			AppId: appId,
		})
	}

	_, err = br.CreateBanners(ctx, bannerBodys)
	if err != nil {
		return err
	}

	return br.db.DeleteBanners(ctx, bannerIds)
}
