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

// Mock repository
type MockLimitRepository struct {
	mock.Mock
}

func (m *MockLimitRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.CustomerWithLimitRaw, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.CustomerWithLimitRaw), args.Get(1).(int64), args.Error(2)
}

func (m *MockLimitRepository) FindByCustomerId(ctx context.Context, id string) (domain.CustomerWithLimit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.CustomerWithLimit), args.Error(1)
}

func (m *MockLimitRepository) FindById(ctx context.Context, id string) (domain.Limit, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Limit), args.Error(1)
}

func (m *MockLimitRepository) Save(ctx context.Context, limit *domain.Limit) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockLimitRepository) Update(ctx context.Context, limit *domain.Limit) error {
	args := m.Called(ctx, limit)
	return args.Error(0)
}

func (m *MockLimitRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestLimitService_Index(t *testing.T) {
	mockRepo := new(MockLimitRepository)
	svc := service.NewLimit(&config.Config{}, mockRepo)

	mockData := []domain.CustomerWithLimitRaw{
		{
			ID:             "1",
			NIK:            "1234567890",
			FullName:       "John Doe",
			LegalName:      "John Doe",
			PlaceOfBirth:   "Jakarta",
			DateOfBirth:    sql.NullTime{Valid: true, Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
			Salary:         5000000,
			KTPPhotoURL:    "ktp.jpg",
			SelfiePhotoURL: "selfie.jpg",
			LimitID:        "limit-1",
			TenorMonths:    sql.NullInt64{Valid: true, Int64: 12},
			LimitAmount:    sql.NullFloat64{Valid: true, Float64: 10000000},
			Status:         sql.NullString{Valid: true, String: "active"},
		},
	}

	mockRepo.On("FindAll", mock.Anything, 10, 0).Return(mockData, int64(1), nil)

	res, total, err := svc.Index(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, res, 1)
	assert.Equal(t, "John Doe", res[0].FullName)
	assert.Equal(t, 1, len(res[0].Limit))
	assert.Equal(t, "active", res[0].Limit[0].Status)
}

func TestLimitService_Show(t *testing.T) {
	mockRepo := new(MockLimitRepository)
	svc := service.NewLimit(&config.Config{}, mockRepo)

	mockCustomer := domain.CustomerWithLimit{
		ID:             "1",
		NIK:            "1234567890",
		FullName:       "John Doe",
		LegalName:      "John Doe",
		PlaceOfBirth:   "Jakarta",
		DateOfBirth:    sql.NullTime{Valid: true, Time: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
		Salary:         5000000,
		KTPPhotoURL:    "ktp.jpg",
		SelfiePhotoURL: "selfie.jpg",
		Limits: []domain.LimitDetail{
			{
				LimitID:     "limit-1",
				TenorMonths: 12,
				LimitAmount: 10000000,
				Status:      "active",
			},
		},
	}

	mockRepo.On("FindByCustomerId", mock.Anything, "1").Return(mockCustomer, nil)

	res, err := svc.Show(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", res.FullName)
	assert.Equal(t, 1, len(res.Limit))
	assert.Equal(t, "active", res.Limit[0].Status)
}

func TestLimitService_Create(t *testing.T) {
	mockRepo := new(MockLimitRepository)
	svc := service.NewLimit(&config.Config{}, mockRepo)

	req := dto.CreateLimitRequest{
		CustomerId:  "1",
		TenorMonths: 12,
		LimitAmount: 10000000,
		Status:      "active",
	}

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
}

func TestLimitService_Update(t *testing.T) {
	mockRepo := new(MockLimitRepository)
	svc := service.NewLimit(&config.Config{}, mockRepo)

	existing := domain.Limit{
		ID: "limit-1",
	}

	req := dto.UpdateLimitRequest{
		ID:          "limit-1",
		CustomerId:  "1",
		TenorMonths: 12,
		LimitAmount: 10000000,
		Status:      "active",
	}

	mockRepo.On("FindById", mock.Anything, "limit-1").Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
}

func TestLimitService_Delete(t *testing.T) {
	mockRepo := new(MockLimitRepository)
	svc := service.NewLimit(&config.Config{}, mockRepo)

	existing := domain.Limit{ID: "limit-1"}

	mockRepo.On("FindById", mock.Anything, "limit-1").Return(existing, nil)
	mockRepo.On("Delete", mock.Anything, "limit-1").Return(nil)

	err := svc.Delete(context.Background(), "limit-1")

	assert.NoError(t, err)
}
