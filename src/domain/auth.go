package domain

import (
	"context"
	"saldri/test_pt_xyz/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
	LoginWeb(ctx context.Context, req dto.AuthRequest) (dto.AuthWebResponse, error)
}
