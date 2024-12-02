package swipe

import (
	"context"
	"fmt"

	"github.com/a-berahman/dating-app/internal/repository"

	"go.uber.org/zap"
)

type SwipeLogic struct {
	swipeRepo repository.SwipeRepository
	matchRepo repository.MatchRepository
	logger    *zap.Logger
}

func NewSwipeLogic(swipeRepo repository.SwipeRepository, matchRepo repository.MatchRepository, logger *zap.Logger) *SwipeLogic {
	return &SwipeLogic{
		swipeRepo: swipeRepo,
		matchRepo: matchRepo,
		logger:    logger,
	}
}

// ProcessSwipe processes a swipe action and checks for matches
func (sl *SwipeLogic) ProcessSwipe(ctx context.Context, userID, targetUserID uint, swipedRight bool) (bool, uint, error) {
	swipe := repository.Swipe{
		UserID:       userID,
		TargetUserID: targetUserID,
		SwipedRight:  swipedRight,
	}
	if err := sl.swipeRepo.AddSwipe(ctx, &swipe); err != nil {
		sl.logger.Error("Failed to add swipe", zap.Error(err))
		return false, 0, fmt.Errorf("failed to add swipe: %w", err)

	}

	if swipedRight {
		return sl.processPotentialMatch(ctx, userID, targetUserID)

	}

	return false, 0, nil
}
func (sl *SwipeLogic) processPotentialMatch(ctx context.Context, userID, targetUserID uint) (bool, uint, error) {
	matched, err := sl.swipeRepo.CheckForMatch(ctx, userID, targetUserID)
	if err != nil {
		sl.logger.Error("Error checking for match", zap.Uint("userID", userID), zap.Uint("targetUserID", targetUserID), zap.Error(err))
		return false, 0, fmt.Errorf("error checking for match: %w", err)
	}
	if matched {
		matchID, err := sl.matchRepo.CreateOrUpdateMatch(ctx, userID, targetUserID)
		if err != nil {
			sl.logger.Error("Error creating or updating match", zap.Error(err))
			return false, 0, fmt.Errorf("error creating or updating match: %w", err)
		}
		return true, matchID, nil
	}
	return false, 0, nil
}
