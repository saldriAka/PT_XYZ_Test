package dto

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthWebResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
