package service

import (
	"context"
	"database/sql"
	"errors"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"sync"
	"time"

	"github.com/google/uuid"
)

type transactionsService struct {
	cnf                    *config.Config
	transactionsRepository domain.TransactionsRepository
}

func NewTransactions(
	cnf *config.Config,
	transactionsRepository domain.TransactionsRepository,
) domain.TransactionsService {
	return &transactionsService{
		cnf:                    cnf,
		transactionsRepository: transactionsRepository,
	}
}

func (s transactionsService) Index(ctx context.Context, page, limit int) ([]dto.TransactionsData, int64, error) {
	offset := (page - 1) * limit

	result, total, err := s.transactionsRepository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var data []dto.TransactionsData
	for _, v := range result {
		data = append(data, dto.TransactionsData{
			ID:             v.ID,
			ContractNumber: v.ContractNumber,
			Channel:        v.Channel,
			OTRAmount:      v.OTRAmount,
			AdminFee:       v.AdminFee,
			Installment:    v.Installment,
			Interest:       v.Interest,
			AssetName:      v.AssetName,
			TenorMonths:    v.TenorMonths,
			Customer: dto.CustomerDTO{
				ID:             v.Customer.ID,
				NIK:            v.Customer.NIK,
				FullName:       v.Customer.FullName,
				LegalName:      v.Customer.LegalName,
				PlaceOfBirth:   v.Customer.PlaceOfBirth,
				DateOfBirth:    v.Customer.DateOfBirth,
				Salary:         v.Customer.Salary,
				KTPPhotoURL:    v.Customer.KTPPhotoURL,
				SelfiePhotoURL: v.Customer.SelfiePhotoURL,
			},
		})
	}

	return data, total, nil
}

func (s transactionsService) CustomerShow(ctx context.Context, id string) ([]dto.Transactions, error) {
	result, err := s.transactionsRepository.FindByCustomerId(ctx, id)
	if err != nil {
		return nil, err
	}

	var data []dto.Transactions
	for _, v := range result {
		data = append(data, dto.Transactions{
			ID:             v.ID,
			CustomerID:     v.CustomerID,
			ContractNumber: v.ContractNumber,
			Channel:        v.Channel,
			OTRAmount:      v.OTRAmount,
			AdminFee:       v.AdminFee,
			Installment:    v.Installment,
			Interest:       v.Interest,
			AssetName:      v.AssetName,
			TenorMonths:    v.TenorMonths,
		})
	}

	return data, nil
}

func (s transactionsService) Show(ctx context.Context, id string) (dto.TransactionsData, error) {
	transaction, err := s.transactionsRepository.FindById(ctx, id)
	if err != nil {
		return dto.TransactionsData{}, err
	}

	if transaction.ID == "" {
		return dto.TransactionsData{}, errors.New("data transaction tidak ditemukan")
	}

	return dto.TransactionsData{
		ID:             transaction.ID,
		ContractNumber: transaction.ContractNumber,
		Channel:        transaction.Channel,
		OTRAmount:      transaction.OTRAmount,
		AdminFee:       transaction.AdminFee,
		Installment:    transaction.Installment,
		Interest:       transaction.Interest,
		AssetName:      transaction.AssetName,
		TenorMonths:    transaction.TenorMonths,
		Customer: dto.CustomerDTO{
			ID:             transaction.Customer.ID,
			NIK:            transaction.Customer.NIK,
			FullName:       transaction.Customer.FullName,
			LegalName:      transaction.Customer.LegalName,
			PlaceOfBirth:   transaction.Customer.PlaceOfBirth,
			DateOfBirth:    transaction.Customer.DateOfBirth,
			Salary:         transaction.Customer.Salary,
			KTPPhotoURL:    transaction.Customer.KTPPhotoURL,
			SelfiePhotoURL: transaction.Customer.SelfiePhotoURL,
		},
	}, nil
}

func (s transactionsService) Create(ctx context.Context, req dto.CreateTransactionsRequest) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(req.Transactions))
	for _, t := range req.Transactions {
		wg.Add(1)
		go func(tx dto.SingleTransactionRequest) {
			defer wg.Done()

			transaction := domain.Transactions{
				ID:             uuid.NewString(),
				CustomerID:     tx.CustomerID,
				ContractNumber: tx.ContractNumber,
				Channel:        tx.Channel,
				OTRAmount:      tx.OTRAmount,
				AdminFee:       tx.AdminFee,
				Installment:    tx.Installment,
				Interest:       tx.Interest,
				AssetName:      tx.AssetName,
				TenorMonths:    tx.TenorMonths,
				CreatedAt: sql.NullTime{
					Valid: true,
					Time:  time.Now(),
				},
			}

			if err := s.transactionsRepository.Save(ctx, &transaction); err != nil {
				errChan <- err
			}
		}(t)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return err
	}

	return nil
}

func (s transactionsService) Update(ctx context.Context, req dto.UpdateTransactionsRequest) error {
	transaction, err := s.transactionsRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}

	if transaction.ID == "" {
		return errors.New("data transaction tidak ditemukan")
	}

	transaction.Channel = req.Channel
	transaction.OTRAmount = req.OTRAmount
	transaction.AdminFee = req.AdminFee
	transaction.Installment = req.Installment
	transaction.Interest = req.Interest
	transaction.AssetName = req.AssetName
	transaction.TenorMonths = req.TenorMonths
	transaction.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	return s.transactionsRepository.Update(ctx, &transaction)
}

func (s transactionsService) Delete(ctx context.Context, id string) error {
	transaction, err := s.transactionsRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if transaction.ID == "" {
		return errors.New("data transaction tidak ditemukan")
	}

	return s.transactionsRepository.Delete(ctx, transaction.ID)
}
