package repository

import (
	"context"

	"github.com/pkg/errors"
)

// AddSwipe logs a swipe action in the database.
func (r *repo) AddSwipe(ctx context.Context, swipe *Swipe) error {
	return r.db.WithContext(ctx).Create(swipe).Error
}

// CheckForMatch checks if two users have swiped right on each other
func (r *repo) CheckForMatch(ctx context.Context, userID, targetUserID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Swipe{}).
		Where("user_id = ? AND target_user_id = ? AND swiped_right = ?", targetUserID, userID, true).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, "failed to check for match")
	}

	return count > 0, nil
}
