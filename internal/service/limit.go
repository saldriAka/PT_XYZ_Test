package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"time"

	"github.com/google/uuid"
)

type limitService struct {
	cnf             *config.Config
	limitRepository domain.LimitRepository
}

func NewLimit(cnf *config.Config, limitRepository domain.LimitRepository) *limitService {
	return &limitService{
		cnf:             cnf,
		limitRepository: limitRepository,
	}
}

func (s *limitService) Index(ctx context.Context) ([]dto.CustomerLimitData, error) {
	result, err := s.limitRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Map keyed by customer ID to aggregate limits
	customerMap := make(map[string]*dto.CustomerLimitData)

	for _, v := range result {
		// Format DOB
		formattedDOB := ""
		if v.DateOfBirth.Valid {
			formattedDOB = v.DateOfBirth.Time.Format("2006-01-02")
		}

		// Build limit detail if valid
		var limit dto.LimitDetail
		if v.TenorMonths.Valid && v.LimitAmount.Valid && v.Status.Valid {
			limit = dto.LimitDetail{
				TenorMonths: int(v.TenorMonths.Int64),
				LimitAmount: v.LimitAmount.Float64,
				Status:      v.Status.String,
			}
		} else {
			// skip limit if not valid
			continue
		}

		// Check if customer already in map
		cust, exists := customerMap[v.ID]
		if !exists {
			// Create new customer entry
			customerMap[v.ID] = &dto.CustomerLimitData{
				CustomersData: dto.CustomersData{
					ID:             v.ID,
					NIK:            v.NIK,
					FullName:       v.FullName,
					LegalName:      v.LegalName,
					PlaceOfBirth:   v.PlaceOfBirth,
					DateOfBirth:    formattedDOB,
					Salary:         v.Salary,
					KTPPhotoURL:    v.KTPPhotoURL,
					SelfiePhotoURL: v.SelfiePhotoURL,
				},
				Limit: []dto.LimitDetail{limit}, // init slice with first limit
			}
		} else {
			// Append limit to existing customer
			cust.Limit = append(cust.Limit, limit)
		}
	}

	// Convert map to slice
	data := make([]dto.CustomerLimitData, 0, len(customerMap))
	for _, v := range customerMap {
		data = append(data, *v)
	}

	return data, nil
}

func (s *limitService) Show(ctx context.Context, id string) (dto.CustomerLimitData, error) {
	customer, err := s.limitRepository.FindByCustomerId(ctx, id)
	if err != nil {
		return dto.CustomerLimitData{}, err
	}

	formattedDOB := ""
	if customer.DateOfBirth.Valid {
		formattedDOB = customer.DateOfBirth.Time.Format("2006-01-02")
	}

	var limits []dto.LimitDetail
	for _, l := range customer.Limits {
		limits = append(limits, dto.LimitDetail{
			TenorMonths: l.TenorMonths,
			LimitAmount: l.LimitAmount,
			Status:      l.Status,
		})
	}

	result := dto.CustomerLimitData{
		CustomersData: dto.CustomersData{
			ID:             customer.ID,
			NIK:            customer.NIK,
			FullName:       customer.FullName,
			LegalName:      customer.LegalName,
			PlaceOfBirth:   customer.PlaceOfBirth,
			DateOfBirth:    formattedDOB,
			Salary:         customer.Salary,
			KTPPhotoURL:    customer.KTPPhotoURL,
			SelfiePhotoURL: customer.SelfiePhotoURL,
		},
		Limit: limits,
	}

	return result, nil
}

func (s *limitService) Create(ctx context.Context, req dto.CreateLimitRequest) error {
	limit := domain.Limit{
		ID:          uuid.NewString(),
		CustomerId:  req.CustomerId,
		TenorMonths: req.TenorMonths,
		LimitAmount: req.LimitAmount,
		Status:      req.Status,
		CreatedAt:   sql.NullTime{Valid: true, Time: time.Now()},
	}

	return s.limitRepository.Save(ctx, &limit)
}

func (s *limitService) Update(ctx context.Context, req dto.UpdateLimitRequest) error {

	existing, err := s.limitRepository.FindById(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("limit tidak ditemukan: %w", err)
	}

	if existing.ID == "" {
		return errors.New("data limit tidak ditemukan")
	}

	// Update field
	existing.CustomerId = req.CustomerId
	existing.TenorMonths = req.TenorMonths
	existing.LimitAmount = req.LimitAmount
	existing.Status = req.Status
	existing.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}

	return s.limitRepository.Update(ctx, &existing)
}

func (s *limitService) Delete(ctx context.Context, id string) error {

	limit, err := s.limitRepository.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("limit tidak ditemukan: %w", err)
	}

	if limit.ID == "" {
		return errors.New("data limit tidak ditemukan")
	}

	return s.limitRepository.Delete(ctx, limit.ID)
}
