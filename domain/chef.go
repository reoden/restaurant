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
	TableChef = "chef"
)

type ChefRepo struct {
	db *gorm.DB
}

func NewChefRepo() *ChefRepo {
	return &ChefRepo{
		db: pkgs.GetDB(),
	}
}

func (ar *ChefRepo) Create(ctx context.Context, app *entity.Chef) (*entity.Chef, error) {
	err := ar.db.Table(TableChef).Create(&app).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *ChefRepo) Get(ctx context.Context, id int64) (*entity.Chef, error) {
	var app entity.Chef
	err := ar.db.Table(TableChef).
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

func (ar *ChefRepo) GetApps(ctx context.Context, ids []int64) ([]*entity.Chef, error) {
	var apps []*entity.Chef

	// 将ids转换为字符串
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = fmt.Sprintf("%d", id)
	}
	idsStr := strings.Join(idStrings, ",")
	err := ar.db.Table(TableChef).
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

func (ar *ChefRepo) GetByUserId(ctx context.Context, userId int64) (*entity.Chef, error) {
	var app entity.Chef
	err := ar.db.Table(TableChef).
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

func (ar *ChefRepo) Update(ctx context.Context, app *entity.Chef) (*entity.Chef, error) {
	err := ar.db.Table(TableChef).
		Where("id = ?", app.Id).
		Where("status != ?", entity.StatusDeleted).
		Updates(app).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *ChefRepo) UpdateField(ctx context.Context, appId int64, params map[string]interface{}) error {
	err := ar.db.Table(TableChef).
		Where("id = ?", appId).
		Where("status != ? AND status != ?", entity.StatusDeleted, entity.StatusSaved).
		Updates(params).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (ar *ChefRepo) Delete(ctx context.Context, appId int64) error {
	err := ar.db.Table(TableChef).
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

func (ar *ChefRepo) List(ctx context.Context) ([]entity.Chef, int64, error) {
	var apps []entity.Chef
	query := ar.db.Table(TableChef).
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

func (br *ChefRepo) GetAll(ctx context.Context) ([]entity.ChefListBody, error) {
	var param []entity.ChefListBody
	err := br.db.Table(TableChef).
		Where("status = ?", entity.StatusAccepted).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
