package handlers

import (
	"github.com/a-berahman/dating-app/internal/handlers/auth"
	"github.com/a-berahman/dating-app/internal/handlers/match"
	"github.com/a-berahman/dating-app/internal/handlers/swipe"
	"github.com/a-berahman/dating-app/internal/handlers/user"
	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserInterface interface {
	CreateFakeUser(c echo.Context) error
}
type AuthInterface interface {
	Login(c echo.Context) error
}
type MatchInterface interface {
	DiscoverMatches(c echo.Context) error
}
type SwipeInterface interface {
	Swipe(c echo.Context) error
}
type Handler struct {
	UserHandler  UserInterface
	AuthHandler  AuthInterface
	MatchHandler MatchInterface
	SwapHadnler  SwipeInterface
}

// New returns a new Handler
func New(l *logic.Logic, logger *zap.Logger) *Handler {
	return &Handler{
		UserHandler:  user.New(l.UserLogic, logger),
		AuthHandler:  auth.New(l.AuthLogic, logger),
		MatchHandler: match.New(l.MatchLogic, logger),
		SwapHadnler:  swipe.New(l.SwipeLogic, logger),
	}
}
