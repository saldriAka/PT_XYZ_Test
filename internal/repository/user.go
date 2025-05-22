package repository

import (
	"context"
	"saldri/test_pt_xyz/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var usr domain.User
	err := u.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&usr).Error

	if err == gorm.ErrRecordNotFound {
		return domain.User{}, nil
	}
	return usr, err
}

func (u userRepository) FindByCustomerEmail(ctx context.Context, email string) (domain.Customers, error) {
	var usr domain.Customers
	err := u.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&usr).Error

	if err == gorm.ErrRecordNotFound {
		return domain.Customers{}, nil
	}
	return usr, err
}
