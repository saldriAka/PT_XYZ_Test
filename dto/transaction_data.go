package dto

type TransactionsData struct {
	ID             string      `json:"id"`
	Customer       CustomerDTO `json:"customer"`
	ContractNumber string      `json:"contract_number"`
	Channel        string      `json:"channel"`
	OTRAmount      float64     `json:"otr_amount"`
	AdminFee       float64     `json:"admin_fee"`
	Installment    float64     `json:"installment"`
	Interest       float64     `json:"interest"`
	AssetName      string      `json:"asset_name"`
	TenorMonths    int         `json:"tenor_months"`
}

type CustomerDTO struct {
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

type CreateTransactionsRequest struct {
	Transactions []SingleTransactionRequest `json:"transactions"`
}

type SingleTransactionRequest struct {
	CustomerID     string  `json:"customer_id" validate:"required"`
	ContractNumber string  `json:"contract_number" validate:"required"`
	Channel        string  `json:"channel" validate:"required"`
	OTRAmount      float64 `json:"otr_amount" validate:"required"`
	AdminFee       float64 `json:"admin_fee" validate:"required"`
	Installment    float64 `json:"installment"`
	Interest       float64 `json:"interest"`
	AssetName      string  `json:"asset_name" validate:"required"`
	TenorMonths    int     `json:"tenor_months" validate:"required"`
}

type UpdateTransactionsRequest struct {
	ID          string  `json:"id" validate:"required"`
	Channel     string  `json:"channel" validate:"required"`
	OTRAmount   float64 `json:"otr_amount" validate:"required"`
	AdminFee    float64 `json:"admin_fee" validate:"required"`
	Installment float64 `json:"installment"`
	Interest    float64 `json:"interest"`
	AssetName   string  `json:"asset_name" validate:"required"`
	TenorMonths int     `json:"tenor_months" validate:"required"`
}
