package domain

import (
	"context"
	"fmt"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"
	"strings"

	"gorm.io/gorm"
)

const (
	TableNews = "news"
)

type NewsRepo struct {
	db *gorm.DB
}

func NewNewsRepo() *NewsRepo {
	return &NewsRepo{
		db: pkgs.GetDB(),
	}
}

func (ar *NewsRepo) Create(ctx context.Context, app *entity.News) (*entity.News, error) {
	err := ar.db.Table(TableNews).Create(&app).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *NewsRepo) Get(ctx context.Context, id int64) (*entity.News, error) {
	var app entity.News
	err := ar.db.Table(TableNews).
		Where("id = ?", id).
		Where("status != ?", entity.StatusDeleted).
		First(&app).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &app, err
}

func (ar *NewsRepo) GetApps(ctx context.Context, ids []int64) ([]*entity.News, error) {
	var apps []*entity.News

	// 将ids转换为字符串
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = fmt.Sprintf("%d", id)
	}
	idsStr := strings.Join(idStrings, ",")
	err := ar.db.Table(TableNews).
		Where("id IN ?", ids).
		Where("status != ?", entity.StatusDeleted).
		Order(fmt.Sprintf("FIELD(id, %s)", idsStr)).
		Find(&apps).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return apps, err
}

func (ar *NewsRepo) GetByUserId(ctx context.Context, userId int64) (*entity.News, error) {
	var app entity.News
	err := ar.db.Table(TableNews).
		Where("user_id = ?", userId).
		Where("status != ?", entity.StatusDeleted).
		First(&app).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &app, err
}

func (ar *NewsRepo) Update(ctx context.Context, app *entity.News) (*entity.News, error) {
	err := ar.db.Table(TableNews).
		Where("id = ?", app.Id).
		Where("status != ?", entity.StatusDeleted).
		Updates(app).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *NewsRepo) UpdateField(ctx context.Context, appId int64, params map[string]interface{}) error {
	err := ar.db.Table(TableNews).
		Where("id = ?", appId).
		Where("status != ? AND status != ?", entity.StatusDeleted, entity.StatusSaved).
		Updates(params).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (ar *NewsRepo) Delete(ctx context.Context, appId int64) error {
	err := ar.db.Table(TableNews).
		Where("id = ?", appId).
		Updates(map[string]interface{}{
			"status": entity.StatusDeleted,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (ar *NewsRepo) List(ctx context.Context) ([]entity.News, int64, error) {
	var apps []entity.News
	query := ar.db.Table(TableNews).
		Where("address != ''").
		Where("status != ?", entity.StatusDeleted)
	err := query.Order("id desc").
		Find(&apps).
		Error

	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	return apps, total, err
}

func (br *NewsRepo) GetAll(ctx context.Context) ([]entity.NewsListBody, error) {
	var param []entity.NewsListBody
	err := br.db.Table(TableNews).
		Where("status = ?", entity.StatusAccepted).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
