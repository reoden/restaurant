package application

import (
	"context"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type GuestShowListService struct {
	db      *domain.GuestShowListRepo
	as      GuestService
	account AccountService
}

func NewGuestShowListService() *GuestShowListService {
	return &GuestShowListService{
		db:      domain.NewGuestShowListRepo(),
		as:      *NewGuestService(),
		account: *NewAccountService(),
	}
}

func (br *GuestShowListService) Create(ctx context.Context, userId int64, param *entity.GuestShowListBody) (*entity.GuestShowListBody, error) {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}
	accountStatus := user.Status
	if accountStatus != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	banner, err := br.db.GetByAppId(ctx, param.AppId)
	if err != nil {
		return nil, err
	}

	if banner != nil {
		return nil, nil
	}

	currentTime := common.NowInLocal()
	param.CreatedAt = currentTime
	param.Status = entity.BannerShowStatus
	return br.db.Create(ctx, param)
}

func (br *GuestShowListService) CreateBanners(ctx context.Context, param []*entity.GuestShowListBody) ([]*entity.GuestShowListBody, error) {
	currentTime := common.NowInLocal()
	for i := range param {
		param[i].CreatedAt = currentTime
		param[i].Status = entity.BannerShowStatus
	}
	return br.db.CreateBanners(ctx, param)
}

func (br *GuestShowListService) Delete(ctx context.Context, userId, id int64) error {
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

func (br *GuestShowListService) List(ctx context.Context) (*entity.GuestShowListResp, error) {
	apps, total, err := br.db.List(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.GuestShowListResp{
		Total: total,
		Apps:  apps,
	}, nil
}

func (br *GuestShowListService) Update(ctx context.Context, userId int64, appIds []int64) error {
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

	bannerBodys := make([]*entity.GuestShowListBody, 0)
	for _, appId := range appIds {
		bannerBodys = append(bannerBodys, &entity.GuestShowListBody{
			AppId: appId,
		})
	}

	_, err = br.CreateBanners(ctx, bannerBodys)
	if err != nil {
		return err
	}

	return br.db.DeleteBanners(ctx, bannerIds)
}
