package dto

type LimitDetail struct {
	LimitID     string  `json:"limit_id" validate:"required,uuid4"`
	TenorMonths int     `json:"tenor_months"`
	LimitAmount float64 `json:"limit_amount"`
	Status      string  `json:"status"`
}

type CustomerLimitData struct {
	CustomersData
	Limit []LimitDetail `json:"limit"`
}

type CreateLimitRequest struct {
	CustomerId  string  `json:"customer_id" validate:"required"`
	TenorMonths int     `json:"tenor_months" validate:"required,oneof=1 2 3 6"`
	LimitAmount float64 `json:"limit_amount" validate:"required,gt=0"`
	Status      string  `json:"status" validate:"required,oneof=available booked"`
}

type UpdateLimitRequest struct {
	ID          string  `json:"id" validate:"required,uuid4"`
	CustomerId  string  `json:"customer_id" validate:"required,uuid4"`
	TenorMonths int     `json:"tenor_months" validate:"required,oneof=1 2 3 6"`
	LimitAmount float64 `json:"limit_amount" validate:"required,gt=0"`
	Status      string  `json:"status" validate:"required,oneof=available booked"`
}
