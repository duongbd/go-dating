package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, email, password, name string, gender constant.UserGender, dateOfBirth time.Time, lat, lng float64) (uint, error) {
	args := m.Called(ctx, email, password, name, gender, dateOfBirth, lat, lng)
	return 0, args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	if usr, ok := args.Get(0).(*repository.User); ok {
		return usr, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*repository.User, error) {
	return nil, nil
}
func (m *MockUserRepository) Authenticate(ctx context.Context, email, password string) (*repository.User, error) {
	return nil, nil
}
func (m *MockUserRepository) FindPotentialMatches(ctx context.Context, userID uint, filters *repository.MatchFilters, lat, lng float64) ([]repository.User, error) {
	return nil, nil
}
func (m *MockUserRepository) AddSwipe(ctx context.Context, swipe *repository.Swipe) error { return nil }
func (m *MockUserRepository) CreateOrUpdateMatch(ctx context.Context, userID, targetUserID uint) (uint, error) {
	return 0, nil
}
func (m *MockUserRepository) CheckForMatch(ctx context.Context, userID, targetUserID uint) (bool, error) {
	return false, nil
}

func TestUserLogic_RegisterUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	logger, _ := zap.NewDevelopment()
	userLogic := NewUserLogic(mockRepo, logger)

	tests := []struct {
		name          string
		email         string
		password      string
		userName      string
		gender        constant.UserGender
		dob           time.Time
		lat           float64
		lng           float64
		setupMock     func()
		expectedError error
	}{
		{
			name:     "successful user registration",
			email:    "test1@example.com",
			password: "passwooord",
			userName: "test username",
			gender:   constant.UserGenderMale,
			dob:      time.Now(),
			lat:      34.0522,
			lng:      -118.2437,
			setupMock: func() {
				mockRepo.On("FindByEmail", ctx, "test1@example.com").Return(nil, nil) // Corrected nil return
				mockRepo.On("Create", ctx, "test1@example.com", "passwooord", "test username", constant.UserGenderMale, mock.AnythingOfType("time.Time"), 34.0522, -118.2437).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "email already in use",
			email:    "test2@example.com",
			password: "passwooord",
			userName: "test username",
			gender:   constant.UserGenderMale,
			dob:      time.Now(),
			lat:      34.0522,
			lng:      -118.2437,
			setupMock: func() {
				mockRepo.On("FindByEmail", ctx, "test2@example.com").Return(&repository.User{}, nil) // Simulate existing user
			},
			expectedError: constant.ErrEmailInUse,
		},
		{
			name:     "repository error on find by email",
			email:    "test3@example.com",
			password: "passsswwooorrrd",
			userName: "test username",
			gender:   constant.UserGenderMale,
			dob:      time.Now(),
			lat:      34.0522,
			lng:      -118.2437,
			setupMock: func() {
				mockRepo.On("FindByEmail", ctx, "test3@example.com").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			_, err := userLogic.RegisterUser(ctx, WithEmail(tc.email), WithPassword(tc.password),
				WithName(tc.userName), WithGender(tc.gender), WithDOB(tc.dob),
				WithLocation(tc.lat, tc.lng))
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
