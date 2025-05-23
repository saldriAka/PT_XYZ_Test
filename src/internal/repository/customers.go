package repository

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"

	"gorm.io/gorm"
)

type customersRepository struct {
	db *gorm.DB
}

func NewCustomers(db *gorm.DB) domain.CustomersRepository {
	return &customersRepository{
		db: db,
	}
}

func (r *customersRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Customers, int64, error) {
	var customers []domain.Customers
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.Customers{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&customers).Error

	return customers, total, err
}

func (r *customersRepository) FindById(ctx context.Context, id string) (domain.Customers, error) {
	var customer domain.Customers
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&customer).Error

	if err == gorm.ErrRecordNotFound {
		return domain.Customers{}, nil // bisa juga return error custom not found
	}
	return customer, err
}

func (r *customersRepository) FindByIds(ctx context.Context, ids []string) ([]domain.Customers, error) {
	var customers []domain.Customers
	if len(ids) == 0 {
		return customers, nil
	}
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Where("deleted_at IS NULL").
		Find(&customers).Error
	return customers, err
}

func (r *customersRepository) Save(ctx context.Context, customer *domain.Customers) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customersRepository) Update(ctx context.Context, customer *domain.Customers) error {
	return r.db.WithContext(ctx).
		Model(&domain.Customers{}).
		Where("id = ? AND deleted_at IS NULL", customer.ID).
		Updates(customer).Error
}

func (r *customersRepository) UpdateAssets(ctx context.Context, customer *domain.Customers) error {
	return r.db.WithContext(ctx).
		Model(&domain.Customers{}).
		Where("id = ? AND deleted_at IS NULL", customer.ID).
		Updates(map[string]interface{}{
			"ktp_photo_url":    customer.KTPPhotoURL,
			"selfie_photo_url": customer.SelfiePhotoURL,
		}).Error
}

func (r *customersRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Customers{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}
