package logic

import (
	"context"

	"go.uber.org/zap"

	"github.com/a-berahman/dating-app/internal/logic/auth"
	"github.com/a-berahman/dating-app/internal/logic/match"
	"github.com/a-berahman/dating-app/internal/logic/swipe"
	"github.com/a-berahman/dating-app/internal/logic/user"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/a-berahman/dating-app/internal/repository"
)

type UserInterface interface {
	RegisterUser(ctx context.Context, opts ...user.UserOption) (uint, error)
}
type MatchInterface interface {
	FindMatches(ctx context.Context, userID uint, opts ...match.MatchOption) ([]model.UserDTO, error)
}
type AuthInterface interface {
	GenerateToken(ctx context.Context, email, password string) (string, error)
}
type SwipeInterface interface {
	ProcessSwipe(ctx context.Context, userID, targetUserID uint, swipedRight bool) (bool, uint, error)
}

type Logic struct {
	UserLogic  UserInterface
	MatchLogic MatchInterface
	AuthLogic  AuthInterface
	SwipeLogic SwipeInterface
}

// New returns a new Logic
func New(repo *repository.Repository, logger *zap.Logger) *Logic {
	return &Logic{
		UserLogic:  user.NewUserLogic(repo.UserRepo, logger),
		MatchLogic: match.NewMatchLogic(repo.MatchRepo, logger),
		AuthLogic:  auth.NewAuthLogic(repo.UserRepo, logger),
		SwipeLogic: swipe.NewSwipeLogic(repo.SwipeRepo, repo.MatchRepo, logger),
	}
}
