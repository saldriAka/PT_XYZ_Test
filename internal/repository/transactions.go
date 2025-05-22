package repository

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"

	"gorm.io/gorm"
)

type transactionsRepository struct {
	db *gorm.DB
}

func NewTransactions(db *gorm.DB) domain.TransactionsRepository {
	return &transactionsRepository{
		db: db,
	}
}

func (r *transactionsRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Transactions, int64, error) {
	var transactions []domain.Transactions
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&domain.Transactions{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Preload("Customer").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionsRepository) FindById(ctx context.Context, id string) (domain.Transactions, error) {
	var transaction domain.Transactions
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&transaction).Error

	if err == gorm.ErrRecordNotFound {
		return domain.Transactions{}, nil
	}
	return transaction, err
}

func (r *transactionsRepository) FindByCustomerId(ctx context.Context, id string) ([]domain.Transactions, error) {
	var transactions []domain.Transactions
	if len(id) == 0 {
		return transactions, nil
	}
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", id).
		Where("deleted_at IS NULL").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionsRepository) Save(ctx context.Context, transaction *domain.Transactions) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionsRepository) Update(ctx context.Context, transaction *domain.Transactions) error {
	return r.db.WithContext(ctx).
		Model(&domain.Transactions{}).
		Where("id = ? AND deleted_at IS NULL", transaction.ID).
		Updates(transaction).Error
}

func (r *transactionsRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Transactions{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now()).Error
}
