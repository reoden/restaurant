package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableBanner = "banner_id"
)

type BannerRepo struct {
	db *gorm.DB
}

func NewBannerRepo() *BannerRepo {
	return &BannerRepo{
		db: pkgs.GetDB(),
	}
}

func (br *BannerRepo) Create(ctx context.Context, bannerBody *entity.BannerBody) (*entity.BannerBody, error) {
	err := br.db.Table(TableBanner).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *BannerRepo) CreateBanners(ctx context.Context, bannerBody []*entity.BannerBody) ([]*entity.BannerBody, error) {
	err := br.db.Table(TableBanner).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *BannerRepo) Delete(ctx context.Context, id int64) error {
	err := br.db.Table(TableBanner).
		Where("app_id = ?", id).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *BannerRepo) DeleteBanners(ctx context.Context, ids []int64) error {
	err := br.db.Table(TableBanner).
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

func (br *BannerRepo) GetByAppId(ctx context.Context, appId int64) (*entity.BannerBody, error) {
	var param entity.BannerBody
	err := br.db.Table(TableBanner).
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

func (br *BannerRepo) List(ctx context.Context) ([]entity.BannerBody, int64, error) {
	var param []entity.BannerBody
	query := br.db.Table(TableBanner).Where("status = ?", entity.BannerShowStatus)
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

func (br *BannerRepo) GetAll(ctx context.Context) ([]entity.BannerBody, error) {
	var param []entity.BannerBody
	err := br.db.Table(TableBanner).
		Where("status = ?", entity.BannerShowStatus).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
