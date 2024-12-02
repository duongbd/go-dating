package auth

import (
	"cmp"
	"context"
	"os"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/model"
	"github.com/a-berahman/dating-app/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AuthLogic struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewAuthLogic(userRepo repository.UserRepository, logger *zap.Logger) *AuthLogic {
	return &AuthLogic{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (al *AuthLogic) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := al.userRepo.Authenticate(ctx, email, password)
	if err != nil {
		al.logger.Error("Failed to authenticate user", zap.String("email", email), zap.Error(err))
		return "", errors.Wrap(err, "authentication failed")
	}

	duration := constant.JWT_DEFAULT_DURATION_VALUE
	if os.Getenv(constant.JWT_CONFIG_DURATION_KEY) != "" {
		d, err := time.ParseDuration(os.Getenv(constant.JWT_CONFIG_DURATION_KEY))
		if err != nil {
			return "", errors.Wrap(err, "failed to parse the duration")
		}
		duration = d
	}

	expirationTime := time.Now().Add(duration)
	claims := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cmp.Or(os.Getenv(constant.JWT_CONFIG_SECRET_KEY), constant.JWT_DEFAULT_SECRET_VALUE)))
	if err != nil {
		al.logger.Error("Failed to sign token", zap.Error(err))
		return "", errors.Wrap(err, "failed to sign the token")
	}

	al.logger.Debug("Token generated", zap.String("token", tokenString), zap.Time("expires", expirationTime), zapcore.Field{Key: "userId", Type: zapcore.Uint64Type, Integer: int64(user.ID)})
	return tokenString, nil
}
