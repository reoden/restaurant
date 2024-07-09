package application

import (
	"context"
	"errors"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type AccountService struct {
	db *domain.AccountRepo
}

func NewAccountService() *AccountService {
	return &AccountService{
		db: domain.NewAccountRepo(),
	}
}

func (as *AccountService) GetAccount(ctx context.Context, id int64) (*entity.Account, error) {
	account, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (as *AccountService) GetAccountByPhone(ctx context.Context, phone string) (*entity.Account, error) {
	account, err := as.db.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (as *AccountService) GetAccountPhone(ctx context.Context, id int64) (string, error) {
	account, err := as.db.Get(ctx, id)
	if err != nil {
		return "", nil
	}

	return account.Phone, nil
}

func (as *AccountService) GetAccountSelf(ctx context.Context, id int64) (map[string]interface{}, error) {
	account, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"phone": account.Phone,
	}

	return res, nil
}

func (as *AccountService) Create(ctx context.Context, account *entity.Account) (interface{}, error) {
	user, err := as.GetAccountByPhone(ctx, account.Phone)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("the phone is already registered")
	}
	currentTime := common.NowInLocal()
	account.CreatedAt = currentTime

	user, err = as.db.Create(ctx, account)
	if err != nil {
		return nil, errors.New("acount create failed")
	}

	token := pkgs.CreateJWTToken(user.Id)
	return map[string]interface{}{
		"access_token": token.AccessToken,
		"expired_at":   token.ExpiredAt,
	}, nil
}

func (as *AccountService) UpdatePhone(ctx context.Context, id int64, password string) error {
	_, err := as.GetAccount(ctx, id)
	if err != nil {
		return pkgs.AppNotFound
	}

	return as.db.UpdatePhone(ctx, id, password)
}
