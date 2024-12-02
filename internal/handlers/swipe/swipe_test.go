package swipe

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockSwipeLogic struct {
	Err     error
	Match   bool
	MatchID uint
}

func (m *MockSwipeLogic) ProcessSwipe(ctx context.Context, userID, targetUserID uint, swipedRight bool) (bool, uint, error) {
	return m.Match, m.MatchID, m.Err
}

func TestSwipe(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupMock      logic.SwipeInterface
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Successful Swipe with Match",
			requestBody: `{"targetUserId": 2, "preference": "YES"}`,
			setupMock: &MockSwipeLogic{
				Match:   true,
				MatchID: 123,
				Err:     nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"results":{"matched":true,"matchID":123}}`,
		},
		{
			name:        "Successful Swipe without Match",
			requestBody: `{"targetUserId": 3, "preference": "YES"}`,
			setupMock: &MockSwipeLogic{
				Match:   false,
				MatchID: 0,
				Err:     nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"results":{"matched":false}}`,
		},
		{
			name:           "Invalid Request Data",
			requestBody:    `{}`,
			setupMock:      &MockSwipeLogic{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Key: 'SwipeRequest.TargetUserID' Error:Field validation for 'TargetUserID' failed on the 'required' tag\nKey: 'SwipeRequest.Preference' Error:Field validation for 'Preference' failed on the 'required' tag"}`,
		},
		{
			name:        "Internal Server Error on Processing Swipe",
			requestBody: `{"targetUserId": 2, "preference": "YES"}`,
			setupMock: &MockSwipeLogic{
				Match:   false,
				MatchID: 0,
				Err:     errors.New("error processing swipe"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"error processing swipe"}`,
		},
	}

	e := echo.New()
	v := validator.New()
	e.Validator = &Validator{validator: v}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/swipe", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.Set("userID", uint(1))
			logger, _ := zap.NewDevelopment()
			h := New(tc.setupMock, logger)

			if assert.NoError(t, h.Swipe(c)) {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
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
