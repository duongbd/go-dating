package match

import (
	"context"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/a-berahman/dating-app/internal/repository"

	"github.com/a-berahman/dating-app/pkg/geo"
	"github.com/a-berahman/dating-app/pkg/utils"
	"go.uber.org/zap"
)

type MatchLogic struct {
	matchRepo repository.MatchRepository
	logger    *zap.Logger
}

func NewMatchLogic(matchRepo repository.MatchRepository, logger *zap.Logger) *MatchLogic {
	return &MatchLogic{
		matchRepo: matchRepo,
		logger:    logger,
	}
}

// MatchOptions holds the options for finding matches
type MatchOptions struct {
	lat, lng, distance float64
	gender             constant.UserGender
	minAge, maxAge     int
}

type MatchOption func(*MatchOptions) // functional options for match finding

// FindMatches finds potential matches for a user based on the given options
func (ml *MatchLogic) FindMatches(ctx context.Context, userID uint, opts ...MatchOption) ([]model.UserDTO, error) {
	options := newMatchOptions(opts...)
	currentYear := time.Now().Year()
	filters := repository.MatchFilters{
		MaxDistance: options.distance,
		Gender:      string(options.gender),
		MinDOB:      time.Date(currentYear-options.maxAge, 1, 1, 0, 0, 0, 0, time.UTC),
		MaxDOB:      time.Date(currentYear-options.minAge+1, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	users, err := ml.matchRepo.FindPotentialMatches(ctx, userID, &filters, options.lat, options.lng)
	if err != nil {
		ml.logger.Error("Failed to find potential matches", zap.Error(err))
		return nil, err
	}

	return ml.processUsers(users)

}
func newMatchOptions(opts ...MatchOption) MatchOptions {
	mo := MatchOptions{}
	for _, opt := range opts {
		opt(&mo)
	}
	return mo
}
func (ml *MatchLogic) processUsers(users []repository.User) ([]model.UserDTO, error) {
	result := make([]model.UserDTO, len(users))
	ch := make(chan int, len(users))
	for i, user := range users {
		go ml.convertToUserDTO(&user, ch, i, result)
	}
	for range users {
		<-ch
	}
	return result, nil
}

func (ml *MatchLogic) convertToUserDTO(user *repository.User, ch chan int, index int, results []model.UserDTO) {
	defer func() { ch <- 1 }()

	point, err := geo.GeoDecodeString(user.Location)
	if err != nil {
		ml.logger.Error("Failed to decode location", zap.String("location", user.Location), zap.Error(err))
		return
	}

	results[index] = model.UserDTO{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Gender:      constant.UserGender(user.Gender),
		DateOfBirth: user.DateOfBirth,
		Age:         utils.CalculateAge(user.DateOfBirth),
		Location: model.Point{
			Lat: point.Y(),
			Lng: point.X(),
		},
	}
}

func WithLocation(lat, lng float64) MatchOption {
	return func(mo *MatchOptions) {
		mo.lat = lat
		mo.lng = lng
	}
}

func WithDistance(distance float64) MatchOption {
	return func(mo *MatchOptions) {
		mo.distance = distance
	}
}

func WithGender(gender constant.UserGender) MatchOption {
	return func(mo *MatchOptions) {
		mo.gender = gender
	}
}

func WithAgeRange(minAge, maxAge int) MatchOption {
	return func(mo *MatchOptions) {
		mo.minAge = minAge
		mo.maxAge = maxAge
	}
}
