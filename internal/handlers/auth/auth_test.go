package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockAuthLogic struct {
	Token string
	Err   error
}

func (m *MockAuthLogic) GenerateToken(ctx context.Context, email, password string) (string, error) {
	return m.Token, m.Err
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectToken    bool
		setupMock      logic.AuthInterface
	}{
		{
			name: "Successful Login",
			requestBody: map[string]string{
				"email":    "user@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
			setupMock: &MockAuthLogic{
				Token: "token",
				Err:   nil,
			},
		},
		{
			name: "Invalid Credentials",
			requestBody: map[string]string{
				"email":    "user@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
			setupMock: &MockAuthLogic{
				Token: "",
				Err:   errors.New("invalid credentials"),
			},
		},
		{
			name: "Invalid Request Body",
			requestBody: map[string]string{
				"email": "test@",
			},
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			setupMock:      &MockAuthLogic{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			logger, _ := zap.NewDevelopment()
			handler := New(tc.setupMock, logger)
			if assert.NoError(t, handler.Login(c)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				if tc.expectToken {
					assert.Contains(t, rec.Body.String(), "token")
				}
			}

		})
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
