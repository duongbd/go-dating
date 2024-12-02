package match

import (
	"net/http"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/internal/logic/match"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/a-berahman/dating-app/pkg/decode"
	"github.com/a-berahman/dating-app/pkg/geo"
	"github.com/a-berahman/dating-app/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MatchHandler struct {
	matchLogic logic.MatchInterface
	logger     *zap.Logger
}

func New(matchLogic logic.MatchInterface, logger *zap.Logger) *MatchHandler {
	return &MatchHandler{
		matchLogic: matchLogic,
		logger:     logger,
	}
}

// DiscoverMatches fetches potential matches for the user.
func (mh *MatchHandler) DiscoverMatches(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	if userID == 0 {
		mh.logger.Warn("Unauthorized access attempt", zap.Uint("userID", userID))
		return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	var req DiscoverRequest
	if err := decode.DecodeAndValidateRequest(c.Request().Context(), &req, &decode.EchoDecoder{C: c}, &decode.EchoValidator{C: c}); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if req.Distance < 4000 {
		req.Distance = constant.DefaultDiscoveryDistance

	}
	// req.Latitude, req.Longitude, req.Distance, constant.UserGender(req.Gender), req.MinAge, req.MaxAge
	// find matches based on the user's preferences
	matches, err := mh.matchLogic.FindMatches(c.Request().Context(), userID,
		match.WithLocation(req.Latitude, req.Longitude), match.WithDistance(req.Distance),
		match.WithGender(constant.UserGender(req.Gender)), match.WithAgeRange(req.MinAge, req.MaxAge))
	if err != nil {
		mh.logger.Error("Failed to find matches", zap.Error(err))
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	results := formatMatches(matches, &req) // format the matches
	return c.JSON(http.StatusOK, results)
}

func formatMatches(matches []model.UserDTO, req *DiscoverRequest) DiscoverResponse {
	results := make([]MatchResult, 0, len(matches))
	for _, match := range matches {
		result := MatchResult{
			ID:             match.ID,
			Name:           match.Name,
			Gender:         match.Gender,
			Age:            match.Age,
			DistanceFromMe: geo.CalculateDistance(req.Latitude, req.Longitude, match.Location.Lat, match.Location.Lng),
		}
		results = append(results, result)
	}
	return DiscoverResponse{Results: results}
}
