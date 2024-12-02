package swipe

// SwipeRequest defines the structure of the request for the Swipe operation
type SwipeRequest struct {
	TargetUserID uint   `json:"targetUserId" validate:"required"`
	Preference   string `json:"preference" validate:"required,oneof=YES NO"`
}

// SwipeResponse defines the structure of the response for the Swipe operation
type SwipeResponse struct {
	Results SwipeResults `json:"results"`
}

// SwipeResults contains the outcome of the swipe attempt
type SwipeResults struct {
	Matched bool `json:"matched"`
	MatchID uint `json:"matchID,omitempty"`
}
