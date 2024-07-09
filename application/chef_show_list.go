package application

import (
	"context"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type ChefShowListService struct {
	db      *domain.ChefShowListRepo
	as      ChefService
	account AccountService
}

func NewChefShowListService() *ChefShowListService {
	return &ChefShowListService{
		db:      domain.NewChefShowListRepo(),
		as:      *NewChefService(),
		account: *NewAccountService(),
	}
}

func (br *ChefShowListService) Create(ctx context.Context, userId int64, param *entity.ChefShowListBody) (*entity.ChefShowListBody, error) {
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

func (br *ChefShowListService) CreateBanners(ctx context.Context, param []*entity.ChefShowListBody) ([]*entity.ChefShowListBody, error) {
	currentTime := common.NowInLocal()
	for i := range param {
		param[i].CreatedAt = currentTime
		param[i].Status = entity.BannerShowStatus
	}
	return br.db.CreateBanners(ctx, param)
}

func (br *ChefShowListService) Delete(ctx context.Context, userId, id int64) error {
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

func (br *ChefShowListService) List(ctx context.Context) (*entity.ChefShowListResp, error) {
	apps, total, err := br.db.List(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.ChefShowListResp{
		Total: total,
		Apps:  apps,
	}, nil
}

func (br *ChefShowListService) Update(ctx context.Context, userId int64, appIds []int64) error {
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

	bannerBodys := make([]*entity.ChefShowListBody, 0)
	for _, appId := range appIds {
		bannerBodys = append(bannerBodys, &entity.ChefShowListBody{
			AppId: appId,
		})
	}

	_, err = br.CreateBanners(ctx, bannerBodys)
	if err != nil {
		return err
	}

	return br.db.DeleteBanners(ctx, bannerIds)
}
