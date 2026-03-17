package models

type LeaderboardEntry struct {
	// UserID   uint    `json:"user_id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}

type LeaderboardTeam struct {
	RoundNumber int     `json:"round_number,omitempty"`
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	GlobalRank  int     `json:"rank"`
}

type LeaderboardResponse struct {
	Entries []LeaderboardEntry `json:"entries"`
	Total   int                `json:"total"`
}

type EventDetail struct {
	Name         string `json:"event_name"`
	TotalCourts  int    `json:"total_courts"`
	GameType     string `json:"game_type"`
	Location     string `json:"location"`
	PlayDate     string `json:"play_date"`
	TotalPlayers int    `json:"total_players"`
	TeamType     string `json:"team_type"`
}

type ScoreEntry struct {
	RoundNumber int     `json:"round_number"`
	TeamName    string  `json:"team_name"`
	PlayerName  string  `json:"name" gorm:"column:player_name"`
	Score       float64 `json:"score"`
	// Rank        int     `json:"rank"`
}

func (LeaderboardEntry) TableName() string {
	return "leaderboard_entries"
}
