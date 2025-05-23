package repository

import (
	"context"
	"saldri/test_pt_xyz/domain"
	"time"

	"gorm.io/gorm"
)

type limitRepository struct {
	db *gorm.DB
}

func NewLimit(db *gorm.DB) domain.LimitRepository {
	return &limitRepository{db: db}
}

func (r *limitRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.CustomerWithLimitRaw, int64, error) {
	var results []domain.CustomerWithLimitRaw
	var total int64

	err := r.db.WithContext(ctx).
		Table("customers").
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Table("customers").
		Select("customers.*, `limit`.tenor_months, `limit`.id as limit_id,`limit`.limit_amount, `limit`.status").
		Joins("LEFT JOIN `limit` ON `limit`.customer_id = customers.id AND `limit`.deleted_at IS NULL").
		Where("customers.deleted_at IS NULL").
		Order("customers.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	return results, total, err
}

func (r *limitRepository) FindByCustomerId(ctx context.Context, id string) (domain.CustomerWithLimit, error) {
	var customer domain.CustomerWithLimit

	// Ambil data customer
	err := r.db.WithContext(ctx).
		Table("customers").
		Select("id, nik, full_name, legal_name, place_of_birth, date_of_birth, salary, ktp_photo_url, selfie_photo_url").
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(&customer).Error
	if err != nil {
		return domain.CustomerWithLimit{}, err
	}

	// Ambil semua limits milik customer
	var limits []domain.LimitDetail
	err = r.db.WithContext(ctx).
		Table("limit").
		Select("id as limit_id, tenor_months, limit_amount, status").
		Where("customer_id = ? AND deleted_at IS NULL", id).
		Scan(&limits).Error
	if err != nil {
		return domain.CustomerWithLimit{}, err
	}

	customer.Limits = limits
	return customer, nil
}

func (r *limitRepository) FindById(ctx context.Context, id string) (domain.Limit, error) {
	var limit domain.Limit

	err := r.db.WithContext(ctx).
		Model(&domain.Limit{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&limit).Error

	if err != nil {
		return domain.Limit{}, err
	}

	return limit, nil
}

func (r *limitRepository) Save(ctx context.Context, limit *domain.Limit) error {

	return r.db.WithContext(ctx).Create(limit).Error
}

func (r *limitRepository) Update(ctx context.Context, limit *domain.Limit) error {
	return r.db.WithContext(ctx).
		Model(&domain.Limit{}).
		Where("id = ? AND deleted_at IS NULL", limit.ID).
		Updates(limit).Error
}

func (r *limitRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Limit{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}
