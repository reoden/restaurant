package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo() *AccountRepo {
	return &AccountRepo{
		db: pkgs.GetDB(),
	}
}

func (ar *AccountRepo) Get(ctx context.Context, id int64) (*entity.Account, error) {
	var account entity.Account
	err := ar.db.Table("account").Where("id = ?", id).First(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &account, nil
}

func (ar *AccountRepo) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	err := ar.db.Table("account").Create(account).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return account, nil
}

func (ar *AccountRepo) GetByPhone(ctx context.Context, phone string) (*entity.Account, error) {
	var account entity.Account
	err := ar.db.Table("account").Where("phone = ?", phone).First(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &account, nil
}

func (ar *AccountRepo) UpdatePhone(ctx context.Context, id int64, phone string) error {
	_, err := ar.Get(ctx, id)
	if err != nil {
		return sferror.WithStack(err)
	}

	err = ar.db.Table("account").Where("id = ?", id).Updates(map[string]interface{}{
		"phone": phone,
	}).Error

	if err != nil {
		return sferror.WithStack(err)
	}
	return nil
}
