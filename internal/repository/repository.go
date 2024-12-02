package repository

import (
	"context"
	"time"

	"github.com/a-berahman/dating-app/constant"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data interaction.
type UserRepository interface {
	Create(ctx context.Context, email, password, name string, gender constant.UserGender, dateOfBirth time.Time, lat, lng float64) (uint, error)
	Authenticate(ctx context.Context, email, password string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// MatchRepository defines the interface for match data interaction.
type MatchRepository interface {
	CreateOrUpdateMatch(ctx context.Context, userID, targetUserID uint) (uint, error)
	FindPotentialMatches(ctx context.Context, userID uint, filters *MatchFilters, lat, lng float64) ([]User, error)
}

// SwipeRepository defines the interface for swipe data interaction.
type SwipeRepository interface {
	AddSwipe(ctx context.Context, swipe *Swipe) error
	CheckForMatch(ctx context.Context, userID, targetUserID uint) (bool, error)
}

// Repository handles the operations with the database
type Repository struct {
	UserRepo  UserRepository
	MatchRepo MatchRepository
	SwipeRepo SwipeRepository
}
type repo struct {
	db *gorm.DB
}

// New creates a new instance of repository layer
func New(db *gorm.DB) *Repository {
	return &Repository{
		UserRepo:  &repo{db: db},
		MatchRepo: &repo{db: db},
		SwipeRepo: &repo{db: db},
	}
}
