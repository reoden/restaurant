package application

import (
	"context"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"time"
)

const (
	expireActivation = 60
)

type VerificationCodeService struct {
	db *domain.VerificationCodeRepo
}

func NewVerificationCodeService() *VerificationCodeService {
	return &VerificationCodeService{
		db: domain.NewVerificationCodeRepo(),
	}
}

func (vs *VerificationCodeService) GetCode(ctx context.Context, phone string) (*entity.VerificationCode, error) {
	verificationCode, err := vs.db.GetCode(ctx, phone)
	if err != nil {
		return nil, err
	}

	if verificationCode == nil {
		return nil, nil
	}

	currentTime := common.NowInLocal()
	if currentTime.After(verificationCode.ExpiredAt) {
		vs.Delete(ctx, verificationCode.Phone, verificationCode.Code)
		return nil, err
	}

	return verificationCode, nil
}

func (vs *VerificationCodeService) Get(ctx context.Context, phone, verificationCode string) (*entity.VerificationCode, error) {
	code, err := vs.db.Get(ctx, phone, verificationCode)
	if err != nil {
		return nil, err
	}

	currentTime := common.NowInLocal()
	if currentTime.After(code.ExpiredAt) {
		vs.Delete(ctx, code.Phone, code.Code)
		return nil, err
	}
	return code, nil
}

// Create 创建
func (vs *VerificationCodeService) Create(ctx context.Context, param entity.VerificationCode) (*entity.VerificationCode, error) {
	currentTime := common.NowInLocal()
	param.CreatedAt = currentTime
	param.UpdatedAt = currentTime
	param.ExpiredAt = currentTime.Add(expireActivation * time.Second)
	param.Status = entity.CodeSended

	return vs.db.Create(ctx, &param)
}

func (vs *VerificationCodeService) Delete(ctx context.Context, phone, verificationCode string) error {
	code, err := vs.db.Get(ctx, phone, verificationCode)
	if err != nil {
		return err
	}

	return vs.db.Delete(ctx, code.Id)
}

func (vs *VerificationCodeService) UpdateStatus(ctx context.Context, phone, verificationCode string) error {
	code, err := vs.Get(ctx, phone, verificationCode)
	if err != nil {
		return err
	}

	return vs.db.UpdateField(ctx, code.Id, map[string]interface{}{
		"status": entity.CodeUsed,
	})
}

func (vs *VerificationCodeService) CountTimes(ctx context.Context, phone string) (int64, error) {
	return vs.db.CountTimes(ctx, phone)
}
