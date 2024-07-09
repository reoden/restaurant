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
	TableApp = "app"
)

type AppRepo struct {
	db *gorm.DB
}

func NewAppRepo() *AppRepo {
	return &AppRepo{
		db: pkgs.GetDB(),
	}
}

func (ar *AppRepo) Create(ctx context.Context, app *entity.Application) (*entity.Application, error) {
	err := ar.db.Table(TableApp).Create(&app).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *AppRepo) Get(ctx context.Context, id int64) (*entity.Application, error) {
	var app entity.Application
	err := ar.db.Table(TableApp).
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

func (ar *AppRepo) GetApps(ctx context.Context, ids []int64) ([]*entity.Application, error) {
	var apps []*entity.Application
	err := ar.db.Table(TableApp).
		Where("id IN ?", ids).
		Where("status != ?", entity.StatusDeleted).
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

func (ar *AppRepo) GetAppsOrder(ctx context.Context, ids []int64) ([]*entity.Application, error) {
	var apps []*entity.Application

	// 将ids转换为字符串
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = fmt.Sprintf("%d", id)
	}
	idsStr := strings.Join(idStrings, ",")
	err := ar.db.Table(TableApp).
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

func (ar *AppRepo) GetByUserId(ctx context.Context, userId int64) (*entity.Application, error) {
	var app entity.Application
	err := ar.db.Table(TableApp).
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

func (ar *AppRepo) Update(ctx context.Context, app *entity.Application) (*entity.Application, error) {
	param := map[string]interface{}{
		"status":        app.Status,
		"phone":         app.Phone,
		"pictures":      app.Pictures,
		"name":          app.Name,
		"describe":      app.Describe,
		"address":       app.Address,
		"post_code":     app.PostCode,
		"post_name":     app.PostName,
		"work_begin_at": app.WorkBeginAt,
		"work_end_at":   app.WorkEndAt,
		"have_vege":     app.HaveVege,
	}
	err := ar.db.Table(TableApp).
		Where("id = ?", app.Id).
		Where("status != ?", entity.StatusDeleted).
		Updates(param).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, err
}

func (ar *AppRepo) UpdateField(ctx context.Context, appId int64, params map[string]interface{}) error {
	err := ar.db.Table(TableApp).
		Where("id = ?", appId).
		Where("status != ? AND status != ?", entity.StatusDeleted, entity.StatusSaved).
		Updates(params).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (ar *AppRepo) Delete(ctx context.Context, appId int64) error {
	err := ar.db.Table(TableApp).
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

func (ar *AppRepo) List(ctx context.Context, offset, limit int) ([]entity.Application, int64, error) {
	var apps []entity.Application
	query := ar.db.Table(TableApp).
		Where("address != ''").
		Where("status != ? AND status != ?", entity.StatusDeleted, entity.StatusSaved)
	err := query.Order("id desc").
		Offset(offset).
		Limit(limit).
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

func (ar *AppRepo) SearchList(ctx context.Context, offset, limit int) ([]entity.Application, int64, error) {
	var apps []entity.Application
	query := ar.db.Table(TableApp).
		Where("address != ''").
		Where("status = ?", entity.StatusAccepted)
	err := query.Order("id desc").
		Offset(offset).
		Limit(limit).
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

func (ar *AppRepo) SearchByName(ctx context.Context, name string) ([]entity.Application, int64, error) {
	var apps []entity.Application
	query := ar.db.Table(TableApp).
		Where("address != ''").
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).
		Where("status = ?", entity.StatusAccepted)
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

func (ar *AppRepo) SearchById(ctx context.Context, id int64) ([]entity.Application, int64, error) {
	var apps []entity.Application
	query := ar.db.Table(TableApp).
		Where("address != ''").
		Where("id LIKE ?", fmt.Sprintf("%%%d%%", id)).
		Where("status = ?", entity.StatusAccepted)
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
