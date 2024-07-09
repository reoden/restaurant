package domain

import (
	"context"
	"restaurant/common"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableCode = "verification_code"
)

type VerificationCodeRepo struct {
	db *gorm.DB
}

func NewVerificationCodeRepo() *VerificationCodeRepo {
	return &VerificationCodeRepo{
		db: pkgs.GetDB(),
	}
}

func (vr *VerificationCodeRepo) GetCode(ctx context.Context, phone string) (*entity.VerificationCode, error) {
	currentTime := common.NowInLocal()
	var verificationCode entity.VerificationCode
	err := vr.db.Table(TableCode).
		Where("phone = ?", phone).
		Where("expired_at > ?", currentTime).
		First(&verificationCode).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &verificationCode, nil
}

func (vr *VerificationCodeRepo) Get(ctx context.Context, phone, code string) (*entity.VerificationCode, error) {
	currentTime := common.NowInLocal()
	var verificationCode entity.VerificationCode
	err := vr.db.Table(TableCode).
		Where("code = ?", code).
		Where("phone = ?", phone).
		Where("status = ?", entity.CodeSended).
		Where("expired_at > ?", currentTime).
		First(&verificationCode).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &verificationCode, nil
}

func (vr *VerificationCodeRepo) Create(ctx context.Context, verificationCode *entity.VerificationCode) (*entity.VerificationCode, error) {
	err := vr.db.Table(TableCode).Create(verificationCode).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return verificationCode, nil
}

func (vr *VerificationCodeRepo) UpdateField(ctx context.Context, codeId int64, params map[string]interface{}) error {
	err := vr.db.Table(TableCode).
		Where("id = ?", codeId).
		Where("status = ?", entity.CodeSended).
		Updates(params).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (vr *VerificationCodeRepo) Delete(ctx context.Context, verificationCodeId int64) error {
	err := vr.db.Table(TableCode).
		Where("id = ?", verificationCodeId).
		Updates(map[string]interface{}{
			"status": entity.CodeExpired,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (vr *VerificationCodeRepo) CountTimes(ctx context.Context, phone string) (int64, error) {
	today := common.NowInLocal()
	ts := common.TimeToDateStr(today)
	query := vr.db.Table(TableCode).Where("phone = ?", phone).Where("DATE(created_at) = ?", ts)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return -1, sferror.WithStack(err)
	}
	return total, nil
}
