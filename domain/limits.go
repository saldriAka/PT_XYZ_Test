package domain

import (
	"context"
	"database/sql"
	"saldri/test_pt_xyz/dto"
)

type Limit struct {
	ID          string       `db:"id"`
	CustomerId  string       `db:"customer_id"`
	TenorMonths int          `db:"tenor_months"`
	LimitAmount float64      `db:"limit_amount"`
	Status      string       `db:"status"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

func (Limit) TableName() string {
	return "limit"
}

type CustomerWithLimitRaw struct {
	ID             string
	NIK            string
	FullName       string
	LegalName      string
	PlaceOfBirth   string
	DateOfBirth    sql.NullTime
	Salary         float64
	KTPPhotoURL    string
	SelfiePhotoURL string
	LimitID        string
	TenorMonths    sql.NullInt64
	LimitAmount    sql.NullFloat64
	Status         sql.NullString
}

type LimitDetail struct {
	LimitID     string  `json:"limit_id" validate:"required,uuid4"`
	TenorMonths int     `json:"tenor_months"`
	LimitAmount float64 `json:"limit_amount"`
	Status      string  `json:"status"`
}

type CustomerWithLimit struct {
	ID             string        `json:"id"`
	NIK            string        `json:"nik"`
	FullName       string        `json:"full_name"`
	LegalName      string        `json:"legal_name"`
	PlaceOfBirth   string        `json:"place_of_birth"`
	DateOfBirth    sql.NullTime  `json:"date_of_birth"`
	Salary         float64       `json:"salary"`
	KTPPhotoURL    string        `json:"ktp_photo_url"`
	SelfiePhotoURL string        `json:"selfie_photo_url"`
	Limits         []LimitDetail `gorm:"-" json:"limit"`
}

type LimitRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]CustomerWithLimitRaw, int64, error)
	FindByCustomerId(ctx context.Context, id string) (CustomerWithLimit, error)
	FindById(ctx context.Context, id string) (Limit, error)
	Save(ctx context.Context, book *Limit) error
	Update(ctx context.Context, book *Limit) error
	Delete(ctx context.Context, id string) error
}

type LimitService interface {
	Index(ctx context.Context, limit, offset int) ([]dto.CustomerLimitData, int64, error)
	Show(ctx context.Context, id string) (dto.CustomerLimitData, error)
	Create(ctx context.Context, req dto.CreateLimitRequest) error
	Update(ctx context.Context, req dto.UpdateLimitRequest) error
	Delete(ctx context.Context, id string) error
}
