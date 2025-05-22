package service_test

import (
	"context"
	"errors"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/service"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepo) FindByCustomerEmail(ctx context.Context, email string) (domain.Customers, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.Customers), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := domain.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	mockRepo := new(mockUserRepo)
	cfg := &config.Config{
		Jwt: config.Jwt{
			Key: "testsecret",
			Exp: 15,
		},
	}
	svc := service.NewAuth(cfg, mockRepo)

	t.Run("success login", func(t *testing.T) {
		mockRepo.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)
		resp, err := svc.Login(context.Background(), dto.AuthRequest{
			Email:    user.Email,
			Password: "correctpassword",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)

		parsedToken, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.Key), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo.On("FindByEmail", mock.Anything, user.Email).Return(user, nil)
		_, err := svc.Login(context.Background(), dto.AuthRequest{
			Email:    user.Email,
			Password: "wrongpassword",
		})
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(domain.User{}, errors.New("not found"))
		_, err := svc.Login(context.Background(), dto.AuthRequest{
			Email:    "notfound@example.com",
			Password: "password",
		})
		assert.Error(t, err)
	})
}

func TestAuthService_LoginWeb(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("webpass"), bcrypt.DefaultCost)
	user := domain.User{
		ID:       "cust-456",
		Email:    "customer@example.com",
		Password: string(hashedPassword),
	}

	mockRepo := new(mockUserRepo)
	cfg := &config.Config{}
	svc := service.NewAuth(cfg, mockRepo)

	t.Run("success login web", func(t *testing.T) {
		mockRepo.On("FindByCustomerEmail", mock.Anything, user.Email).Return(user, nil)
		resp, err := svc.LoginWeb(context.Background(), dto.AuthRequest{
			Email:    user.Email,
			Password: "webpass",
		})
		assert.NoError(t, err)
		assert.Equal(t, user.ID, resp.ID)
		assert.Equal(t, user.Email, resp.Email)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo.On("FindByCustomerEmail", mock.Anything, user.Email).Return(user, nil)
		_, err := svc.LoginWeb(context.Background(), dto.AuthRequest{
			Email:    user.Email,
			Password: "wrongpass",
		})
		assert.Error(t, err)
		assert.Equal(t, "password salah", err.Error())
	})

	t.Run("email not found", func(t *testing.T) {
		mockRepo.On("FindByCustomerEmail", mock.Anything, "notfound@web.com").Return(domain.User{}, errors.New("not found"))
		_, err := svc.LoginWeb(context.Background(), dto.AuthRequest{
			Email:    "notfound@web.com",
			Password: "any",
		})
		assert.Error(t, err)
		assert.Equal(t, "email tidak ditemukan", err.Error())
	})
}
