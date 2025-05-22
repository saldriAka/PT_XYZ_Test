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

type customersService struct {
	cnf                 *config.Config
	customersRepository domain.CustomersRepository
}

func NewCustomers(
	cnf *config.Config,
	customersRepository domain.CustomersRepository,
) domain.CustomersService {
	return &customersService{
		cnf:                 cnf,
		customersRepository: customersRepository,
	}
}

func (s customersService) Index(ctx context.Context, page, limit int) ([]dto.CustomersData, int64, error) {
	offset := (page - 1) * limit

	result, total, err := s.customersRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var data []dto.CustomersData
	for _, v := range result {
		var formattedDOB string
		if v.DateOfBirth.Valid {
			formattedDOB = v.DateOfBirth.Time.Format("2006-01-02")
		} else {
			formattedDOB = "-"
		}

		data = append(data, dto.CustomersData{
			ID:             v.ID,
			NIK:            v.NIK,
			FullName:       v.FullName,
			LegalName:      v.LegalName,
			PlaceOfBirth:   v.PlaceOfBirth,
			DateOfBirth:    formattedDOB,
			KTPPhotoURL:    v.KTPPhotoURL,
			SelfiePhotoURL: v.SelfiePhotoURL,
			Salary:         v.Salary,
		})
	}

	return data, total, nil
}

func (s customersService) Show(ctx context.Context, id string) (dto.CustomersShowData, error) {
	customer, err := s.customersRepository.FindById(ctx, id)
	if err != nil {
		return dto.CustomersShowData{}, err
	}

	if customer.ID == "" {
		return dto.CustomersShowData{}, errors.New("data customer tidak ditemukan")
	}

	var formattedDOB string
	if customer.DateOfBirth.Valid {
		formattedDOB = customer.DateOfBirth.Time.Format("2006-01-02")
	} else {
		formattedDOB = "-"
	}

	return dto.CustomersShowData{
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
	}, nil

}

func (s customersService) Create(ctx context.Context, req dto.CreateCustomersRequest) error {

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date_of_birth: %w", err)
	}

	customer := domain.Customers{
		ID:           uuid.NewString(),
		NIK:          req.NIK,
		FullName:     req.FullName,
		LegalName:    req.LegalName,
		PlaceOfBirth: req.PlaceOfBirth,
		DateOfBirth:  sql.NullTime{Valid: true, Time: dob},
		Salary:       req.Salary,
		CreatedAt: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	return s.customersRepository.Save(ctx, &customer)
}

func (s customersService) Update(ctx context.Context, req dto.UpdateCustomersRequest) error {
	customer, err := s.customersRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}

	if customer.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date_of_birth: %w", err)
	}

	customer.FullName = req.FullName
	customer.NIK = req.NIK
	customer.LegalName = req.LegalName
	customer.PlaceOfBirth = req.PlaceOfBirth
	customer.DateOfBirth = sql.NullTime{Valid: true, Time: dob}
	customer.Salary = req.Salary
	customer.KTPPhotoURL = req.KTPPhotoURL
	customer.SelfiePhotoURL = req.SelfiePhotoURL
	customer.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}

	return s.customersRepository.Update(ctx, &customer)
}

func (s customersService) UpdateAssets(ctx context.Context, id string, req dto.UpdateAssetsCustomersRequest) error {
	customer := &domain.Customers{
		ID:             id,
		KTPPhotoURL:    req.KTPPhotoURL,
		SelfiePhotoURL: req.SelfiePhotoURL,
	}
	return s.customersRepository.UpdateAssets(ctx, customer)
}

func (s customersService) Delete(ctx context.Context, id string) error {
	customer, err := s.customersRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if customer.ID == "" {
		return errors.New("data customer tidak ditemukan")
	}

	return s.customersRepository.Delete(ctx, customer.ID)
}
