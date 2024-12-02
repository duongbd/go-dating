package swipe

import (
	"context"
	"errors"
	"testing"

	"github.com/a-berahman/dating-app/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockMatchRepository struct {
	mock.Mock
}

type MockSwipeRepository struct {
	mock.Mock
}

func (m *MockSwipeRepository) AddSwipe(ctx context.Context, swipe *repository.Swipe) error {
	args := m.Called(ctx, swipe)
	return args.Error(0)
}

func (m *MockSwipeRepository) CheckForMatch(ctx context.Context, userID, targetUserID uint) (bool, error) {
	args := m.Called(ctx, userID, targetUserID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMatchRepository) CreateOrUpdateMatch(ctx context.Context, userID, targetUserID uint) (uint, error) {
	args := m.Called(ctx, userID, targetUserID)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockMatchRepository) FindPotentialMatches(ctx context.Context, userID uint, filters *repository.MatchFilters, lat, lng float64) ([]repository.User, error) {
	return nil, nil
}

func TestSwipeLogic_ProcessSwipe(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		targetUserID    uint
		swipedRight     bool
		setupSwipeMock  func(m *MockSwipeRepository)
		setupMatchMock  func(m *MockMatchRepository)
		expectedMatch   bool
		expectedMatchID uint
		expectedErr     error
	}{
		{
			name:         "successful swipe with match",
			userID:       1,
			targetUserID: 2,
			swipedRight:  true,
			setupSwipeMock: func(m *MockSwipeRepository) {
				m.On("AddSwipe", mock.Anything, mock.AnythingOfType("*repository.Swipe")).Return(nil)
				m.On("CheckForMatch", mock.Anything, uint(1), uint(2)).Return(true, nil)
			},
			setupMatchMock: func(m *MockMatchRepository) {
				m.On("CreateOrUpdateMatch", mock.Anything, uint(1), uint(2)).Return(uint(100), nil)
			},
			expectedMatch:   true,
			expectedMatchID: 100,
			expectedErr:     nil,
		},
		{
			name:         "Successful swipe without match",
			userID:       1,
			targetUserID: 2,
			swipedRight:  true,
			setupSwipeMock: func(m *MockSwipeRepository) {
				m.On("AddSwipe", mock.Anything, mock.AnythingOfType("*repository.Swipe")).Return(nil)
				m.On("CheckForMatch", mock.Anything, uint(1), uint(2)).Return(false, nil)
			},
			setupMatchMock: func(m *MockMatchRepository) {
				m.On("CreateOrUpdateMatch", mock.Anything, uint(1), uint(2)).Return(uint(100), nil)
			},
			expectedMatch:   false,
			expectedMatchID: 0,
			expectedErr:     nil,
		},
		{
			name:         "failed to add sipe",
			userID:       1,
			targetUserID: 2,
			swipedRight:  true,
			setupSwipeMock: func(m *MockSwipeRepository) {
				m.On("AddSwipe", mock.Anything, mock.AnythingOfType("*repository.Swipe")).Return(errors.New("database error"))
			},
			setupMatchMock: func(m *MockMatchRepository) {
				m.On("CreateOrUpdateMatch", mock.Anything, uint(1), uint(2)).Return(uint(100), nil)
			},
			expectedMatch:   false,
			expectedMatchID: 0,
			expectedErr:     errors.New("failed to add swipe: database error"),
		},
	}

	logger, _ := zap.NewDevelopment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSwipeRepo := new(MockSwipeRepository)
			mockMatchRepo := new(MockMatchRepository)
			tt.setupMatchMock(mockMatchRepo)
			tt.setupSwipeMock(mockSwipeRepo)
			logic := NewSwipeLogic(mockSwipeRepo, mockMatchRepo, logger)

			matched, matchID, err := logic.ProcessSwipe(context.Background(), tt.userID, tt.targetUserID, tt.swipedRight)

			assert.Equal(t, tt.expectedMatch, matched)
			if matched {
				assert.Equal(t, tt.expectedMatchID, matchID)
			}
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
