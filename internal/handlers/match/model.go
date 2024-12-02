package match

import "github.com/a-berahman/dating-app/constant"

// DiscoverRequest represents the request to discover potential matches
type DiscoverRequest struct {
	Latitude  float64 `query:"lat" validate:"required"`
	Longitude float64 `query:"lng" validate:"required"`
	Distance  float64 `query:"distance" validate:"gte=0"`
	MinAge    int     `query:"minAge"`
	MaxAge    int     `query:"maxAge"`
	Gender    string  `query:"gender" validate:"required,gender"`
}

// MatchResult represents a single potential match
type MatchResult struct {
	ID             uint                `json:"id"`
	Name           string              `json:"name"`
	Gender         constant.UserGender `json:"gender"`
	Age            int                 `json:"age"`
	DistanceFromMe int                 `json:"distanceFromMe"`
}

// DiscoverResponse represents the collection of match results
type DiscoverResponse struct {
	Results []MatchResult `json:"results"`
}
