package user

import (
	"context"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/repository"
	"go.uber.org/zap"
)

// UserLogic handles business logic for user operations
type UserLogic struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// UserOptions holds the options for user registration
type UserOptions struct {
	email    string
	password string
	name     string
	gender   constant.UserGender
	dob      time.Time
	lat, lng float64
}

type UserOption func(*UserOptions) // functional options for user registration

// NewUserLogic creates a new instance of UserLogic
func NewUserLogic(repo repository.UserRepository, logger *zap.Logger) *UserLogic {
	return &UserLogic{
		userRepo: repo,
		logger:   logger,
	}
}

// RegisterUser registers a new user
func (ul *UserLogic) RegisterUser(ctx context.Context, opts ...UserOption) (uint, error) {
	options := newUserOptions(opts...)
	user, err := ul.userRepo.FindByEmail(ctx, options.email)
	if err != nil {
		ul.logger.Error("failed to find user by email", zap.Error(err))
		return 0, err
	}
	if user != nil {
		ul.logger.Info("email already in use", zap.String("email", options.email))
		return 0, constant.ErrEmailInUse
	}

	userId, err := ul.userRepo.Create(ctx, options.email, options.password, options.name, options.gender, options.dob, options.lat, options.lng)
	if err != nil {
		ul.logger.Error("failed to create user", zap.Error(err))
		return 0, err
	}

	ul.logger.Info("user registered successfully", zap.Uint("userId", userId))
	return userId, nil
}

func newUserOptions(opts ...UserOption) UserOptions {
	uo := UserOptions{}
	for _, opt := range opts {
		opt(&uo)
	}
	return uo
}

func WithEmail(email string) UserOption {
	return func(uo *UserOptions) {
		uo.email = email
	}
}

func WithPassword(password string) UserOption {
	return func(uo *UserOptions) {
		uo.password = password
	}
}

func WithName(name string) UserOption {
	return func(uo *UserOptions) {
		uo.name = name
	}
}

func WithGender(gender constant.UserGender) UserOption {
	return func(uo *UserOptions) {
		uo.gender = gender
	}
}

func WithDOB(dob time.Time) UserOption {
	return func(uo *UserOptions) {
		uo.dob = dob
	}
}

func WithLocation(lat, lng float64) UserOption {
	return func(uo *UserOptions) {
		uo.lat = lat
		uo.lng = lng
	}
}
