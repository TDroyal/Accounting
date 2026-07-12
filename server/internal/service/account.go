package service

import (
	"context"

	"github.com/TDroyal/Accounting/server/internal/model"
	"github.com/TDroyal/Accounting/server/internal/repository"
	"gorm.io/gorm"
)

type AccountService struct {
	repo repository.AccountRepo
}

func NewAccountService(r repository.AccountRepo) *AccountService {
	return &AccountService{repo: r}
}

func (s *AccountService) List(ctx context.Context, userID uint64) ([]model.Account, error) {
	return s.repo.List(ctx, userID)
}

func (s *AccountService) Create(ctx context.Context, userID uint64, a *model.Account) error {
	a.UserID = userID
	if a.Currency == "" {
		a.Currency = "CNY"
	}
	return s.repo.Create(ctx, a)
}

func (s *AccountService) Update(ctx context.Context, userID, id uint64, fields map[string]interface{}) error {
	if err := s.repo.Update(ctx, userID, id, fields); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *AccountService) Delete(ctx context.Context, userID, id uint64) error {
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return nil
}
