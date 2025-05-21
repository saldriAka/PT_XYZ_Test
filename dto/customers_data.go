package dto

type CustomersData struct {
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

type CustomersShowData struct {
	CustomersData
}

type CreateCustomersRequest struct {
	NIK            string  `json:"nik" validate:"required"`
	FullName       string  `json:"full_name" validate:"required"`
	LegalName      string  `json:"legal_name" validate:"required"`
	PlaceOfBirth   string  `json:"place_of_birth" validate:"required"`
	DateOfBirth    string  `json:"date_of_birth" validate:"required"`
	Salary         float64 `json:"salary" validate:"required,gt=0"`
	KTPPhotoURL    string  `json:"ktp_photo_url" validate:"required,url"`
	SelfiePhotoURL string  `json:"selfie_photo_url" validate:"required,url"`
}

type UpdateCustomersRequest struct {
	ID             string  `json:"id" validate:"required,uuid4"`
	FullName       string  `json:"full_name" validate:"required"`
	LegalName      string  `json:"legal_name" validate:"required"`
	PlaceOfBirth   string  `json:"place_of_birth" validate:"required"`
	DateOfBirth    string  `json:"date_of_birth" validate:"required"`
	Salary         float64 `json:"salary" validate:"required,gt=0"`
	KTPPhotoURL    string  `json:"ktp_photo_url" validate:"required,url"`
	SelfiePhotoURL string  `json:"selfie_photo_url" validate:"required,url"`
}
