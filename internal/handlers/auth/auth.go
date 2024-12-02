package auth

import (
	"net/http"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/pkg/decode"
	"github.com/a-berahman/dating-app/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AuthHandler is a handler for authentication operations
type AuthHandler struct {
	authLogic logic.AuthInterface
	logger    *zap.Logger
}

// New creates a new handler for authentication operations
func New(auth logic.AuthInterface, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authLogic: auth,
		logger:    logger,
	}
}

func (ah *AuthHandler) Login(c echo.Context) error {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var req LoginRequest
	if err := decode.DecodeAndValidateRequest(c.Request().Context(), &req, &decode.EchoDecoder{C: c}, &decode.EchoValidator{C: c}); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := ah.authLogic.GenerateToken(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		ah.logger.Error("Invalid credentials", zap.String("email", req.Email), zap.Error(err))
		return utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
	}

	ah.logger.Info("User logged in", zap.String("email", req.Email))
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
