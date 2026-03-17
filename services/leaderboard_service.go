package services

import (
	// "fmt"

	"scoring_app/models"
	"scoring_app/repositories"
)

type LeaderboardService interface {
	GetLeaderboardAllPlayers(page, pageSize int) ([]models.LeaderboardEntry, int64, error)
	GetLeaderboardByMatchID(matchID string) ([]models.LeaderboardTeam, int, error)
	GetEventDetail(matchID string) (*models.EventDetail, error)
	GetScoresByMatchID(matchID string) ([]models.ScoreEntry, error)
	// GetLeaderboardBySession(session string, page, pageSize int) ([]models.LeaderboardEntry, int, error)
	// GetOverallLeaderboard(page, pageSize int) ([]models.LeaderboardEntry, int, error)
}

type LeaderboardServiceImpl struct {
	repo repositories.LeaderboardRepository
}

func NewLeaderboardService(repo repositories.LeaderboardRepository) LeaderboardService {
	return &LeaderboardServiceImpl{repo: repo}
}

func (s *LeaderboardServiceImpl) GetLeaderboardAllPlayers(page, pageSize int) ([]models.LeaderboardEntry, int64, error) {
	return s.repo.GetLeaderboardAllPlayers(page, pageSize)
}

func (s *LeaderboardServiceImpl) GetLeaderboardByMatchID(matchID string) ([]models.LeaderboardTeam, int, error) {
	return s.repo.GetLeaderboardByMatchID(matchID)
}

func (s *LeaderboardServiceImpl) GetEventDetail(matchID string) (*models.EventDetail, error) {
	return s.repo.GetEventDetail(matchID)
}

func (s *LeaderboardServiceImpl) GetScoresByMatchID(matchID string) ([]models.ScoreEntry, error) {
	return s.repo.GetScoresByMatchID(matchID)
}

// func (s *LeaderboardService) GetLeaderboardBySession(session string, page, pageSize int) ([]models.LeaderboardEntry, int, error) {
// 	return s.repo.GetLeaderboardBySession(session, page, pageSize)
// }

// func (s *LeaderboardService) GetOverallLeaderboard(page, pageSize int) ([]models.LeaderboardEntry, int, error) {
// 	return s.repo.GetOverallLeaderboard(page, pageSize)
// }
