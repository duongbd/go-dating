package repository

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateOrUpdateMatch creates or updates a match entry between two users
func (r *repo) CreateOrUpdateMatch(ctx context.Context, userID, targetUserID uint) (uint, error) {

	var match Match
	err := r.db.WithContext(ctx).Where("user_id = ? AND target_user_id = ?", userID, targetUserID).
		Or("user_id = ? AND target_user_id = ?", targetUserID, userID).First(&match).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.Wrap(err, "checking for existing match failed")
	}

	if match.ID == 0 {
		match = Match{UserID: userID, TargetUserID: targetUserID}
		if err := r.db.WithContext(ctx).Create(&match).Error; err != nil {
			return 0, errors.Wrap(err, "creating match failed")
		}
	}

	return match.ID, nil
}

// FindPotentialMatches finds other users excluding the given user and their swipes and applying filters if provided
func (r *repo) FindPotentialMatches(ctx context.Context, userID uint, filters *MatchFilters, lat, lng float64) ([]User, error) {
	var users []User
	subQuery := r.db.WithContext(ctx).Select("target_user_id").Where("user_id = ?", userID).Table("swipes")
	query := r.db.Model(&User{}).
		Where("users.id <> ?", userID).
		Not("users.id IN (?)", subQuery).
		Where("date_of_birth >= ?", filters.MinDOB).
		Where("date_of_birth <= ?", filters.MaxDOB)
	if filters.Gender != "" {
		query = query.Where("gender = ?", filters.Gender)
	}

	if filters != nil && filters.MaxDistance > 0 && lat != 0 && lng != 0 {
		query = query.Where("ST_DWithin(location::geography, ST_MakePoint(?, ?)::geography, ?)",
			lng, lat, filters.MaxDistance)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
