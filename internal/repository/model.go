package repository

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	gorm.Model
	Email       string
	Password    string
	Name        string
	Gender      string
	DateOfBirth time.Time
	Location    string
}

// MatchFilters represents the filters that can be applied when searching for matches
type MatchFilters struct {
	Gender      string
	MinDOB      time.Time
	MaxDOB      time.Time
	MaxDistance float64
}

// Swipe represents the swipe action taken by a user on another user's profile
type Swipe struct {
	gorm.Model
	UserID       uint
	TargetUserID uint
	SwipedRight  bool
}

// Match represents a mutual like between two users
type Match struct {
	gorm.Model
	UserID       uint
	TargetUserID uint
	Matched      bool
}
