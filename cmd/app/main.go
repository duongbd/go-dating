package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/a-berahman/dating-app/internal/handlers"
	"github.com/a-berahman/dating-app/internal/logic"
	"github.com/a-berahman/dating-app/internal/repository"

	customMiddleware "github.com/a-berahman/dating-app/pkg/middleware"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func main() {

	logger := setupLogger()
	defer logger.Sync()

	e := setupEcho()
	db := setupDatabase(e)
	setupRoutes(e, setupHandlers(db, logger), logger)

	startHTTPServer(e)

}
func setupLogger() *zap.Logger {
	var logger *zap.Logger
	var err error
	if cmp.Or(os.Getenv("LOG_ENV"), "development") == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	return logger
}
func setupEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	v.RegisterValidation("gender", genderValidation)

	e.Validator = &Validator{validator: v}
	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	return e
}
func setupDatabase(e *echo.Echo) *gorm.DB {

	db, err := gorm.Open(postgres.Open(concatConnectionString()), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("Error connecting to database: ", err)
	}
	if cmp.Or(os.Getenv("MIGRATION_ENBABLED"), "TRUE") == "TRUE" {
		if err := db.AutoMigrate(&repository.User{}, &repository.Match{}, &repository.Swipe{}); err != nil {
			e.Logger.Fatal("Error auto-migrating database: ", err)
		}
	}
	return db
}

func concatConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cmp.Or(os.Getenv("DB_HOST"), "localhost"),
		cmp.Or(os.Getenv("DB_USER"), "user"),
		cmp.Or(os.Getenv("DB_PASS"), "password"),
		cmp.Or(os.Getenv("DB_NAME"), "datingapp"), cmp.Or(os.Getenv("DB_PORT"), "5432"))
}
func setupHandlers(db *gorm.DB, logger *zap.Logger) *handlers.Handler {
	userRepository := repository.New(db)
	userLogic := logic.New(userRepository, logger)
	return handlers.New(userLogic, logger)
}
func setupRoutes(e *echo.Echo, handler *handlers.Handler, logger *zap.Logger) {
	// e.POST("api/v1/user", handler.UserHandler.RegisterUser)
	// Path of the routs are defined based on the problem statement
	e.POST("/user/create", handler.UserHandler.CreateFakeUser)
	e.POST("/login", handler.AuthHandler.Login)

	e.POST("/swipe", handler.SwapHadnler.Swipe, customMiddleware.UserAuthMiddleware(logger))
	e.GET("/discover", handler.MatchHandler.DiscoverMatches, customMiddleware.UserAuthMiddleware(logger))

}
func startHTTPServer(e *echo.Echo) {
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", cmp.Or(os.Getenv("PORT"), "8080"))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
func genderValidation(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == string(constant.UserGenderMale) || gender == string(constant.UserGenderFemale)
}
