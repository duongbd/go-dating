package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	User *repository.User
	Err  error
}

func (m *MockUserRepository) Authenticate(ctx context.Context, email, password string) (*repository.User, error) {
	return m.User, m.Err
}

func (m *MockUserRepository) Create(ctx context.Context, email, password, name string, gender constant.UserGender, dateOfBirth time.Time, lat, lng float64) (uint, error) {
	return 0, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	return nil, nil
}
func TestGenerateToken(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tests := []struct {
		name          string
		setupMock     repository.UserRepository
		email         string
		password      string
		expectedToken string
		expectError   bool
	}{
		{
			name: "successful token generation",
			setupMock: &MockUserRepository{
				User: &repository.User{
					Model: gorm.Model{
						ID: 1,
					},
				},
				Err: nil,
			},

			email:       "test@example.com",
			password:    "password",
			expectError: false,
		},
		{
			name: "authentication failure",
			setupMock: &MockUserRepository{
				User: &repository.User{},
				Err:  errors.New("authentication failed"),
			},

			email:       "test@example.com",
			password:    "password",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			authLogic := NewAuthLogic(tt.setupMock, logger)

			token, err := authLogic.GenerateToken(context.Background(), tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token, "Token should not be empty")
			}
		})
	}
}
