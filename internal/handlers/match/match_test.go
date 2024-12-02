package match

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/internal/logic/match"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MockMatchLogic is a mock type for the MatchLogic interface
type MockMatchLogic struct {
	Users   []model.UserDTO
	Err     error
	Match   bool
	MatchID uint
}

func (m *MockMatchLogic) FindMatches(ctx context.Context, userID uint, opts ...match.MatchOption) ([]model.UserDTO, error) {
	return m.Users, m.Err
}
func (m *MockMatchLogic) ProcessSwipe(ctx context.Context, userID, targetUserID uint, swipedRight bool) (bool, uint, error) {
	return m.Match, m.MatchID, m.Err
}

func TestDiscoverHandler(t *testing.T) {
	e := echo.New()
	v := validator.New()
	genderValdiation := func(fl validator.FieldLevel) bool {
		gender := fl.Field().String()
		return gender == "MALE" || gender == "FEMALE"
	}
	v.RegisterValidation("gender", genderValdiation)
	e.Validator = &Validator{validator: v}

	scenarios := []struct {
		name           string
		requestPath    string
		setupMock      logic.MatchInterface
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Successful Discovery",
			requestPath: "/discover?lat=34.0522&lng=-118.2437&distance=10&minAge=18&maxAge=30&gender=MALE",
			setupMock: &MockMatchLogic{
				Users: []model.UserDTO{
					{
						ID:     1,
						Name:   "test name",
						Gender: constant.UserGenderMale,
						Age:    25,
						Location: model.Point{
							Lat: 40.7128,
							Lng: -74.0060,
						},
					},
				},
				Err: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"results":[{"id":1,"name":"test name","gender":"MALE","age":25,"distanceFromMe":3935}]}`,
		},
		{
			name:           "Invalid Parameters",
			requestPath:    "/discover?lat=34.0522&lng=-118.2437&distance=-10&minAge=30&maxAge=18&gender=MALE",
			setupMock:      &MockMatchLogic{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Key: 'DiscoverRequest.Distance' Error:Field validation for 'Distance' failed on the 'gte' tag"}`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, scenario.requestPath, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userID", uint(1))
			logger, _ := zap.NewDevelopment()
			h := New(scenario.setupMock, logger)

			if assert.NoError(t, h.DiscoverMatches(c)) {
				assert.Equal(t, scenario.expectedStatus, rec.Code)
				assert.JSONEq(t, scenario.expectedBody, rec.Body.String())
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
