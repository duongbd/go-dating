package repository

import (
	"context"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/pkg/geo"
	hash "github.com/a-berahman/dating-app/pkg/hash"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Create saves a new user in the database
func (r *repo) Create(ctx context.Context, email, password, name string, gender constant.UserGender, dateOfBirth time.Time, lat, lng float64) (uint, error) {
	hashedPassword, err := hash.Generate([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.Wrap(err, "hashing password")
	}

	location, err := geo.GeoEncode(lat, lng)
	if err != nil {
		return 0, errors.Wrap(err, "encoding location")
	}

	user := &User{
		Email:       email,
		Password:    string(hashedPassword),
		Name:        name,
		Location:    location,
		Gender:      string(gender),
		DateOfBirth: dateOfBirth,
	}

	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// FindByEmail finds a user by email
func (r *repo) FindByEmail(ctx context.Context, email string) (*User, error) {

	var user User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// Authenticate checks if credentials are correct
func (r *repo) Authenticate(ctx context.Context, email, password string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.Wrap(err, "finding user")
	}
	if !hash.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}
