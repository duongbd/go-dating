package user

import (
	"net/http"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/internal/logic/user"
	"github.com/a-berahman/dating-app/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type UserHandler struct {
	userLogic logic.UserInterface
	logger    *zap.Logger
}

// New creates a new handler for user operations
func New(logic logic.UserInterface, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userLogic: logic,
		logger:    logger,
	}
}

// CreateFakeUser creates a fake user and returns the user's information
func (h *UserHandler) CreateFakeUser(c echo.Context) error {

	fakeUser := utils.GenerateFakeUser()
	userID, err := h.userLogic.RegisterUser(c.Request().Context(),
		user.WithEmail(fakeUser.Email),
		user.WithPassword(fakeUser.Password),
		user.WithName(fakeUser.Name),
		user.WithGender(fakeUser.Gender),
		user.WithDOB(fakeUser.DateOfBirth),
		user.WithLocation(fakeUser.Latitude, fakeUser.Longitude))
	if err != nil {
		h.logger.Error("Failed to create fake user", zap.Error(err))
		return utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create fake user")
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"result": echo.Map{
			"id":       userID,
			"email":    fakeUser.Email,
			"password": fakeUser.Password,
			"name":     fakeUser.Name,
			"gender":   fakeUser.Gender,
			"age":      utils.CalculateAge(fakeUser.DateOfBirth),
		},
	})
}
