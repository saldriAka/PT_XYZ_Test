package domain

import (
	"context"
	"database/sql"
	"saldri/test_pt_xyz/dto"
)

type Transactions struct {
	ID             string       `db:"id"`
	CustomerID     string       `db:"customer_id"`
	Customer       Customer     `gorm:"foreignKey:CustomerID"`
	ContractNumber string       `db:"contract_number"`
	Channel        string       `db:"channel"`
	OTRAmount      float64      `db:"otr_amount"`
	AdminFee       float64      `db:"admin_fee"`
	Installment    float64      `db:"installment"`
	Interest       float64      `db:"interest"`
	AssetName      string       `db:"asset_name"`
	TenorMonths    int          `db:"tenor_months"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	DeletedAt      sql.NullTime `db:"deleted_at"`
}

type Customer struct {
	ID             string  `json:"id"`
	NIK            string  `json:"nik"`
	FullName       string  `json:"full_name"`
	LegalName      string  `json:"legal_name"`
	PlaceOfBirth   string  `json:"place_of_birth"`
	DateOfBirth    string  `json:"date_of_birth"`
	Salary         float64 `json:"salary"`
	KTPPhotoURL    string  `json:"ktp_photo_url"`
	SelfiePhotoURL string  `json:"selfie_photo_url"`
}

type TransactionsRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]Transactions, int64, error)
	FindById(ctx context.Context, id string) (Transactions, error)
	FindByCustomerId(ctx context.Context, id string) ([]Transactions, error)
	Save(ctx context.Context, book *Transactions) error
	Update(ctx context.Context, book *Transactions) error
	Delete(ctx context.Context, id string) error
}

type TransactionsService interface {
	Index(ctx context.Context, limit, offset int) ([]dto.TransactionsData, int64, error)
	Show(ctx context.Context, id string) (dto.TransactionsData, error)
	CustomerShow(ctx context.Context, id string) ([]dto.Transactions, error)
	Create(ctx context.Context, req dto.CreateTransactionsRequest) error
	Update(ctx context.Context, req dto.UpdateTransactionsRequest) error
	Delete(ctx context.Context, id string) error
}
