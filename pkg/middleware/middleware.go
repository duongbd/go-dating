package middleware

import (
	"cmp"
	"net/http"
	"os"
	"strings"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// UserAuthMiddleware authenticates users and sets their ID in the context
func UserAuthMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				logger.Debug("No Authorization header provided")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header is required"})
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				logger.Debug("Authorization header does not start with Bearer")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header must start with Bearer"})
			}

			tokenStr := authHeader[len(bearerPrefix):]
			claims, err := parseToken(tokenStr)
			if err != nil {
				logger.Error("Token parsing failed", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
			}

			c.Set("userID", claims.UserID)
			return next(c)
		}
	}
}

func parseToken(tokenStr string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cmp.Or(os.Getenv(constant.JWT_CONFIG_SECRET_KEY), constant.JWT_DEFAULT_SECRET_VALUE)), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
