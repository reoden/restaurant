package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableNewsShowList = "news_show_list"
)

type NewsShowListRepo struct {
	db *gorm.DB
}

func NewNewsShowListRepo() *NewsShowListRepo {
	return &NewsShowListRepo{
		db: pkgs.GetDB(),
	}
}

func (br *NewsShowListRepo) Create(ctx context.Context, bannerBody *entity.NewsShowListBody) (*entity.NewsShowListBody, error) {
	err := br.db.Table(TableNewsShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *NewsShowListRepo) CreateBanners(ctx context.Context, bannerBody []*entity.NewsShowListBody) ([]*entity.NewsShowListBody, error) {
	err := br.db.Table(TableNewsShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *NewsShowListRepo) Delete(ctx context.Context, id int64) error {
	err := br.db.Table(TableNewsShowList).
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

func (br *NewsShowListRepo) DeleteBanners(ctx context.Context, ids []int64) error {
	err := br.db.Table(TableNewsShowList).
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

func (br *NewsShowListRepo) GetByAppId(ctx context.Context, appId int64) (*entity.NewsShowListBody, error) {
	var param entity.NewsShowListBody
	err := br.db.Table(TableNewsShowList).
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

func (br *NewsShowListRepo) List(ctx context.Context) ([]entity.NewsShowListBody, int64, error) {
	var param []entity.NewsShowListBody
	query := br.db.Table(TableNewsShowList).Where("status = ?", entity.BannerShowStatus)
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

func (br *NewsShowListRepo) GetAll(ctx context.Context) ([]entity.NewsShowListBody, error) {
	var param []entity.NewsShowListBody
	err := br.db.Table(TableNewsShowList).
		Where("status = ?", entity.BannerShowStatus).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
