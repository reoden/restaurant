package application

import (
	"context"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type NewsShowListService struct {
	db      *domain.NewsShowListRepo
	as      NewsService
	account AccountService
}

func NewNewsShowListService() *NewsShowListService {
	return &NewsShowListService{
		db:      domain.NewNewsShowListRepo(),
		as:      *NewNewsService(),
		account: *NewAccountService(),
	}
}

func (br *NewsShowListService) Create(ctx context.Context, userId int64, param *entity.NewsShowListBody) (*entity.NewsShowListBody, error) {
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

func (br *NewsShowListService) CreateBanners(ctx context.Context, param []*entity.NewsShowListBody) ([]*entity.NewsShowListBody, error) {
	currentTime := common.NowInLocal()
	for i := range param {
		param[i].CreatedAt = currentTime
		param[i].Status = entity.BannerShowStatus
	}
	return br.db.CreateBanners(ctx, param)
}

func (br *NewsShowListService) Delete(ctx context.Context, userId, id int64) error {
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

func (br *NewsShowListService) List(ctx context.Context) (*entity.NewsShowListResp, error) {
	apps, total, err := br.db.List(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.NewsShowListResp{
		Total: total,
		Apps:  apps,
	}, nil
}

func (br *NewsShowListService) Update(ctx context.Context, userId int64, appIds []int64) error {
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

	bannerBodys := make([]*entity.NewsShowListBody, 0)
	for _, appId := range appIds {
		bannerBodys = append(bannerBodys, &entity.NewsShowListBody{
			AppId: appId,
		})
	}

	_, err = br.CreateBanners(ctx, bannerBodys)
	if err != nil {
		return err
	}

	return br.db.DeleteBanners(ctx, bannerIds)
}
