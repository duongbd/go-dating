package utils

import (
	"time"

	"github.com/a-berahman/dating-app/constant"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/labstack/echo/v4"
)

type FakeUserInfo struct {
	Email       string
	Password    string
	Name        string
	Gender      constant.UserGender
	DateOfBirth time.Time
	Latitude    float64
	Longitude   float64
}

// GenerateFakeUser generates fake user data.
func GenerateFakeUser() FakeUserInfo {
	var gender constant.UserGender
	if gofakeit.Bool() {
		gender = constant.UserGenderMale
	} else {
		gender = constant.UserGenderFemale
	}

	return FakeUserInfo{
		Email:       gofakeit.Email(),
		Name:        gofakeit.Name(),
		Password:    gofakeit.Password(true, true, true, true, false, 12),
		Gender:      gender,
		DateOfBirth: gofakeit.DateRange(time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)),
		Latitude:    gofakeit.Latitude(),
		Longitude:   gofakeit.Longitude(),
	}
}

// CalculateAge calculates age given the date of birth
func CalculateAge(dob time.Time) int {
	age := time.Now().Year() - dob.Year()
	if time.Now().Before(time.Date(time.Now().Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.UTC)) {
		age--
	}
	return age
}

// ErrorResponse is a helper to send uniform error responses
func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, echo.Map{"error": message})
}

func GetUserIDFromContext(c echo.Context) uint {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return 0
	}
	return userID
}
