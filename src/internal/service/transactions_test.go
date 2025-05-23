package service_test

import (
	"context"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionsRepo struct {
	mock.Mock
}

func (m *MockTransactionsRepo) FindAll(ctx context.Context, limit, offset int) ([]domain.Transactions, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Transactions), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionsRepo) FindByCustomerId(ctx context.Context, id string) ([]domain.Transactions, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]domain.Transactions), args.Error(1)
}

func (m *MockTransactionsRepo) FindById(ctx context.Context, id string) (domain.Transactions, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Transactions), args.Error(1)
}

func (m *MockTransactionsRepo) Save(ctx context.Context, t *domain.Transactions) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTransactionsRepo) Update(ctx context.Context, t *domain.Transactions) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTransactionsRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTransactionsService_Index(t *testing.T) {
	mockRepo := new(MockTransactionsRepo)
	cfg := &config.Config{}
	svc := service.NewTransactions(cfg, mockRepo)

	dummyTransactions := []domain.Transactions{
		{
			ID:             "trx-1",
			ContractNumber: "CTR-001",
			Channel:        "Online",
			OTRAmount:      100_000_000,
			AdminFee:       1_000_000,
			Installment:    2_000_000,
			Interest:       3.5,
			AssetName:      "Toyota Avanza",
			TenorMonths:    36,
			Customer: domain.Customer{
				ID:       "cust-1",
				NIK:      "1234567890123456",
				FullName: "John Doe",
				Salary:   10_000_000,
			},
		},
	}

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return(dummyTransactions, int64(1), nil)

	data, total, err := svc.Index(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, data, 1)
	assert.Equal(t, "trx-1", data[0].ID)
}

func TestTransactionsService_Create(t *testing.T) {
	mockRepo := new(MockTransactionsRepo)
	cfg := &config.Config{}
	svc := service.NewTransactions(cfg, mockRepo)

	req := dto.CreateTransactionsRequest{
		Transactions: []dto.SingleTransactionRequest{
			{
				CustomerID:     "cust-1",
				ContractNumber: "CTR-001",
				Channel:        "Online",
				OTRAmount:      100_000_000,
				AdminFee:       1_000_000,
				Installment:    2_000_000,
				Interest:       3.5,
				AssetName:      "Avanza",
				TenorMonths:    36,
			},
		},
	}

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Transactions")).Return(nil)

	err := svc.Create(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransactionsService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockTransactionsRepo)
	cfg := &config.Config{}
	svc := service.NewTransactions(cfg, mockRepo)

	mockRepo.On("FindById", mock.Anything, "not-exist").Return(domain.Transactions{}, nil)

	err := svc.Delete(context.Background(), "not-exist")
	assert.EqualError(t, err, "data transaction tidak ditemukan")
}
