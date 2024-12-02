package match

import (
	"context"
	"testing"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/a-berahman/dating-app/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockMatchRepository struct {
	User []repository.User
	Err  error
}

func (m *MockMatchRepository) FindPotentialMatches(ctx context.Context, userID uint, filters *repository.MatchFilters, lat, lng float64) ([]repository.User, error) {
	return m.User, m.Err
}

func (m *MockMatchRepository) CreateOrUpdateMatch(ctx context.Context, userID, targetUserID uint) (uint, error) {
	return 0, nil
}

func TestMatchLogic_FindMatches(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		userID      uint
		lat         float64
		lng         float64
		distance    float64
		gender      constant.UserGender
		minAge      int
		maxAge      int
		mockSetup   repository.MatchRepository
		expected    []model.UserDTO
		expectError bool
	}{

		{
			name:     "successful match retrieval",
			lat:      34.0522,
			lng:      -118.2437,
			distance: 5000,
			gender:   constant.UserGenderMale,
			minAge:   18,
			maxAge:   35,
			mockSetup: &MockMatchRepository{
				User: []repository.User{
					{
						Email:       "test@example.com",
						Name:        "test name",
						Gender:      "FEMALE",
						Location:    "0101000020E610000072D68656DD5E40C08FC2F5285CD44740",
						DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: []model.UserDTO{
				{
					Email:  "test@example.com",
					Name:   "test name",
					Gender: constant.UserGenderFemale,
					Location: model.Point{
						Lat: 47.6590625,
						Lng: -32.74112969955321,
					},
					Age:         34,
					DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ml := NewMatchLogic(tc.mockSetup, logger)
			// tc.lat, tc.lng, tc.distance, tc.gender, tc.minAge, tc.maxAge
			results, err := ml.FindMatches(context.Background(), tc.userID, WithLocation(tc.lat, tc.lng),
				WithDistance(tc.distance), WithGender(tc.gender), WithAgeRange(tc.minAge, tc.maxAge))

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, results)
			}

		})
	}
}
