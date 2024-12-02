package swipe

import (
	"net/http"

	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/pkg/decode"
	"github.com/a-berahman/dating-app/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type SwipeHandler struct {
	swipeLogic logic.SwipeInterface
	logger     *zap.Logger
}

// New creates a new handler for swap operations
func New(swapLogic logic.SwipeInterface, logger *zap.Logger) *SwipeHandler {
	return &SwipeHandler{
		swipeLogic: swapLogic,
		logger:     logger,
	}
}

// Swipe processes a swipe action
func (sh *SwipeHandler) Swipe(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		sh.logger.Debug("Unauthorized swipe attempt")
		return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	var req SwipeRequest
	if err := decode.DecodeAndValidateRequest(c.Request().Context(), &req, &decode.EchoDecoder{C: c}, &decode.EchoValidator{C: c}); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// process the swipe
	matched, matchID, err := sh.swipeLogic.ProcessSwipe(c.Request().Context(), userID, req.TargetUserID, req.Preference == "YES")
	if err != nil {
		sh.logger.Error("Failed to process swipe", zap.Error(err))
		return utils.ErrorResponse(c, http.StatusInternalServerError, "error processing swipe")
	}

	return c.JSON(http.StatusOK, formatSwipeResponse(matched, matchID))
}
func formatSwipeResponse(matched bool, matchID uint) SwipeResponse {

	response := SwipeResponse{
		Results: SwipeResults{
			Matched: matched,
		},
	}

	if matched {
		response.Results.MatchID = matchID
	}

	return response
}
