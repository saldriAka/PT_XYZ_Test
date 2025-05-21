package domain

import (
	"context"
	"database/sql"
	"saldri/test_pt_xyz/dto"
)

type Customers struct {
	ID             string       `db:"id"`
	NIK            string       `db:"nik"`
	FullName       string       `db:"full_name"`
	LegalName      string       `db:"legal_name"`
	PlaceOfBirth   string       `db:"place_of_birth"`
	DateOfBirth    sql.NullTime `db:"date_of_birth"`
	Salary         float64      `db:"salary"`
	KTPPhotoURL    string       `db:"ktp_photo_url"`
	SelfiePhotoURL string       `db:"selfie_photo_url"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	DeletedAt      sql.NullTime `db:"deleted_at"`
}

type CustomersRepository interface {
	FindAll(ctx context.Context) ([]Customers, error)
	FindById(ctx context.Context, id string) (Customers, error)
	FindByIds(ctx context.Context, id []string) ([]Customers, error)
	Save(ctx context.Context, book *Customers) error
	Update(ctx context.Context, book *Customers) error
	Delete(ctx context.Context, id string) error
}

type CustomersService interface {
	Index(ctx context.Context) ([]dto.CustomersData, error)
	Show(ctx context.Context, id string) (dto.CustomersShowData, error)
	Create(ctx context.Context, req dto.CreateCustomersRequest) error
	Update(ctx context.Context, req dto.UpdateCustomersRequest) error
	Delete(ctx context.Context, id string) error
}
