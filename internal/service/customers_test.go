package service_test

import (
	"context"
	"database/sql"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/service"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomersRepository implements domain.CustomersRepository
type MockCustomersRepository struct {
	mock.Mock
}

func (m *MockCustomersRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Customers, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Customers), int64(args.Int(1)), args.Error(2)
}

func (m *MockCustomersRepository) FindById(ctx context.Context, id string) (domain.Customers, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Customers), args.Error(1)
}

func (m *MockCustomersRepository) Save(ctx context.Context, customer *domain.Customers) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomersRepository) Update(ctx context.Context, customer *domain.Customers) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomersRepository) UpdateAssets(ctx context.Context, customer *domain.Customers) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomersRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomersRepository) FindByIds(ctx context.Context, ids []string) ([]domain.Customers, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]domain.Customers), args.Error(1)
}

func TestCustomersService_Index(t *testing.T) {
	mockRepo := new(MockCustomersRepository)
	svc := service.NewCustomers(&config.Config{}, mockRepo)

	expected := []domain.Customers{
		{
			ID:          "123",
			NIK:         "321",
			FullName:    "John Doe",
			LegalName:   "Jonathan Doe",
			DateOfBirth: sql.NullTime{Valid: true, Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
			Salary:      1000000,
		},
	}

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return(expected, 1, nil)

	data, total, err := svc.Index(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "John Doe", data[0].FullName)
	assert.Equal(t, "1990-01-01", data[0].DateOfBirth)
}

func TestCustomersService_Show_NotFound(t *testing.T) {
	mockRepo := new(MockCustomersRepository)
	svc := service.NewCustomers(&config.Config{}, mockRepo)

	mockRepo.On("FindById", mock.Anything, "unknown-id").Return(domain.Customers{}, nil)

	_, err := svc.Show(context.Background(), "unknown-id")
	assert.EqualError(t, err, "data customer tidak ditemukan")
}

func TestCustomersService_Create_InvalidDOB(t *testing.T) {
	mockRepo := new(MockCustomersRepository)
	svc := service.NewCustomers(&config.Config{}, mockRepo)

	req := dto.CreateCustomersRequest{
		DateOfBirth: "invalid-date",
	}

	err := svc.Create(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date_of_birth")
}

func TestCustomersService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockCustomersRepository)
	svc := service.NewCustomers(&config.Config{}, mockRepo)

	mockRepo.On("FindById", mock.Anything, "id-1").Return(domain.Customers{}, nil)

	err := svc.Update(context.Background(), dto.UpdateCustomersRequest{ID: "id-1"})
	assert.EqualError(t, err, "data customer tidak ditemukan")
}

func TestCustomersService_Delete_Success(t *testing.T) {
	mockRepo := new(MockCustomersRepository)
	svc := service.NewCustomers(&config.Config{}, mockRepo)

	mockRepo.On("FindById", mock.Anything, "id-1").Return(domain.Customers{ID: "id-1"}, nil)
	mockRepo.On("Delete", mock.Anything, "id-1").Return(nil)

	err := svc.Delete(context.Background(), "id-1")
	assert.NoError(t, err)
}
