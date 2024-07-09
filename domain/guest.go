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
	TableGuest = "guest"
)

type GuestRepo struct {
	db *gorm.DB
}

func NewGuestRepo() *GuestRepo {
	return &GuestRepo{
		db: pkgs.GetDB(),
	}
}

func (ar *GuestRepo) Create(ctx context.Context, app *entity.Guest) (*entity.Guest, error) {
	err := ar.db.Table(TableGuest).Create(&app).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *GuestRepo) Get(ctx context.Context, id int64) (*entity.Guest, error) {
	var app entity.Guest
	err := ar.db.Table(TableGuest).
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

func (ar *GuestRepo) GetApps(ctx context.Context, ids []int64) ([]*entity.Guest, error) {
	var apps []*entity.Guest

	// 将ids转换为字符串
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = fmt.Sprintf("%d", id)
	}
	idsStr := strings.Join(idStrings, ",")
	err := ar.db.Table(TableGuest).
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

func (ar *GuestRepo) GetByUserId(ctx context.Context, userId int64) (*entity.Guest, error) {
	var app entity.Guest
	err := ar.db.Table(TableGuest).
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

func (ar *GuestRepo) Update(ctx context.Context, app *entity.Guest) (*entity.Guest, error) {
	err := ar.db.Table(TableGuest).
		Where("id = ?", app.Id).
		Where("status != ?", entity.StatusDeleted).
		Updates(app).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *GuestRepo) UpdateField(ctx context.Context, appId int64, params map[string]interface{}) error {
	err := ar.db.Table(TableGuest).
		Where("id = ?", appId).
		Where("status != ? AND status != ?", entity.StatusDeleted, entity.StatusSaved).
		Updates(params).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (ar *GuestRepo) Delete(ctx context.Context, appId int64) error {
	err := ar.db.Table(TableGuest).
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

func (ar *GuestRepo) List(ctx context.Context) ([]entity.Guest, int64, error) {
	var apps []entity.Guest
	query := ar.db.Table(TableGuest).
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

func (br *GuestRepo) GetAll(ctx context.Context) ([]entity.GuestListBody, error) {
	var param []entity.GuestListBody
	err := br.db.Table(TableGuest).
		Where("status = ?", entity.StatusAccepted).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}
