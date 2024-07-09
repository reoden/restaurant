package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableChefShowList = "chef_show_list"
)

type ChefShowListRepo struct {
	db *gorm.DB
}

func NewChefShowListRepo() *ChefShowListRepo {
	return &ChefShowListRepo{
		db: pkgs.GetDB(),
	}
}

func (br *ChefShowListRepo) Create(ctx context.Context, bannerBody *entity.ChefShowListBody) (*entity.ChefShowListBody, error) {
	err := br.db.Table(TableChefShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *ChefShowListRepo) CreateBanners(ctx context.Context, bannerBody []*entity.ChefShowListBody) ([]*entity.ChefShowListBody, error) {
	err := br.db.Table(TableChefShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *ChefShowListRepo) Delete(ctx context.Context, id int64) error {
	err := br.db.Table(TableChefShowList).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *ChefShowListRepo) DeleteBanners(ctx context.Context, ids []int64) error {
	err := br.db.Table(TableChefShowList).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *ChefShowListRepo) GetByAppId(ctx context.Context, appId int64) (*entity.ChefShowListBody, error) {
	var param entity.ChefShowListBody
	err := br.db.Table(TableChefShowList).
		Where("app_id = ?", appId).
		Where("status = ?", entity.BannerShowStatus).
		First(&param).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &param, nil
}

func (br *ChefShowListRepo) List(ctx context.Context) ([]entity.ChefShowListBody, int64, error) {
	var param []entity.ChefShowListBody
	query := br.db.Table(TableChefShowList).Where("status = ?", entity.BannerShowStatus)
	err := query.Find(&param).Error

	if err != nil {
		return param, 0, sferror.WithStack(err)
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return param, 0, sferror.WithStack(err)
	}

	return param, total, nil
}

func (br *ChefShowListRepo) GetAll(ctx context.Context) ([]entity.ChefShowListBody, error) {
	var param []entity.ChefShowListBody
	err := br.db.Table(TableChefShowList).
		Where("status = ?", entity.BannerShowStatus).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
