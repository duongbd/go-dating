package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/internal/logic/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockUserLogic struct {
	UserID uint
	Err    error
}

func (m *MockUserLogic) RegisterUser(ctx context.Context, opts ...user.UserOption) (uint, error) {

	return m.UserID, m.Err
}

func TestRegisterUser(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name           string
		setupMock      logic.UserInterface
		expectedStatus int
	}{
		{
			name:           "Successful User Registration",
			setupMock:      &MockUserLogic{UserID: 1, Err: nil},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failed User Registration",
			setupMock:      &MockUserLogic{UserID: 1, Err: errors.New("error")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserLogic := tt.setupMock
			logger, _ := zap.NewDevelopment()
			handler := New(mockUserLogic, logger)

			req := httptest.NewRequest(http.MethodPost, "/user/create", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if assert.NoError(t, handler.CreateFakeUser(c)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}
