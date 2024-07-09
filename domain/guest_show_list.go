package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableGuestShowList = "guest_show_list"
)

type GuestShowListRepo struct {
	db *gorm.DB
}

func NewGuestShowListRepo() *GuestShowListRepo {
	return &GuestShowListRepo{
		db: pkgs.GetDB(),
	}
}

func (br *GuestShowListRepo) Create(ctx context.Context, bannerBody *entity.GuestShowListBody) (*entity.GuestShowListBody, error) {
	err := br.db.Table(TableGuestShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *GuestShowListRepo) CreateBanners(ctx context.Context, bannerBody []*entity.GuestShowListBody) ([]*entity.GuestShowListBody, error) {
	err := br.db.Table(TableGuestShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *GuestShowListRepo) Delete(ctx context.Context, id int64) error {
	err := br.db.Table(TableGuestShowList).
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

func (br *GuestShowListRepo) DeleteBanners(ctx context.Context, ids []int64) error {
	err := br.db.Table(TableGuestShowList).
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

func (br *GuestShowListRepo) GetByAppId(ctx context.Context, appId int64) (*entity.GuestShowListBody, error) {
	var param entity.GuestShowListBody
	err := br.db.Table(TableGuestShowList).
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

func (br *GuestShowListRepo) List(ctx context.Context) ([]entity.GuestShowListBody, int64, error) {
	var param []entity.GuestShowListBody
	query := br.db.Table(TableGuestShowList).Where("status = ?", entity.BannerShowStatus)
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

func (br *GuestShowListRepo) GetAll(ctx context.Context) ([]entity.GuestShowListBody, error) {
	var param []entity.GuestShowListBody
	err := br.db.Table(TableGuestShowList).
		Where("status = ?", entity.BannerShowStatus).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
