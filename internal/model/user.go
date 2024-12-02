package model

import (
	"time"

	"github.com/a-berahman/dating-app/constant"
)

// User is the model for the user data transfer
type UserDTO struct {
	ID          uint
	Email       string
	Name        string
	Gender      constant.UserGender
	DateOfBirth time.Time
	Location    Point
	Age         int
}

type Point struct {
	Lat float64
	Lng float64
}
