package validators

// MatchEventRequest represents the match event create/update request payload
type BannerEvents struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Location    string `json:"location" binding:"max=255"`
	Image       string `json:"image" binding:"max=65535"`
	Description string `json:"description" binding:"max=65535"`
	StartDate   string `json:"start_date" binding:"required"` // Format: YYYY-MM-DD or YYYY-MM-DD HH:mm:ss
	EndDate     string `json:"end_date" binding:"required"`   // Format: YYYY-MM-DD or YYYY-MM-DD HH:mm:ss
}

// MatchEventResponse represents the match event response
type BannerEventsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
