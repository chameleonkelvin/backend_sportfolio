package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type LeaderboardRepository interface {
	GetLeaderboardAllPlayers(page, pageSize int) ([]models.LeaderboardEntry, int64, error)
	GetLeaderboardByMatchID(matchID string) ([]models.LeaderboardTeam, int, error)
	GetEventDetail(matchID string) (*models.EventDetail, error)
	GetScoresByMatchID(matchID string) ([]models.ScoreEntry, error)
	// GetLeaderboardBySession(session string, page, pageSize int) ([]models.LeaderboardEntry, int, error)
	// GetOverallLeaderboard(page, pageSize int) ([]models.LeaderboardEntry, int, error)
}

type leaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) LeaderboardRepository {
	return &leaderboardRepository{db: db}
}

func (r *leaderboardRepository) GetLeaderboardAllPlayers(page, pageSize int) ([]models.LeaderboardEntry, int64, error) {
	var entries []models.LeaderboardEntry
	var total int64

	if err := r.db.Model(&models.MatchPlayer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Table("match_players").
		Select(`name,
            temp_score as score,
            ROW_NUMBER() OVER (ORDER BY temp_score DESC, name ASC) as rank`).
		Order("temp_score DESC, name ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// GET LEADERBOARD BY MATCH ID
func (r *leaderboardRepository) GetLeaderboardByMatchID(matchID string) ([]models.LeaderboardTeam, int, error) {
	var entries []models.LeaderboardTeam

	query := `
        WITH player_scores AS (
            SELECT mr.team_a_player1_id AS player_id, p1.name AS player_name, mr.score_a AS score
            FROM match_rounds mr
            JOIN match_players p1 ON p1.id = mr.team_a_player1_id
            WHERE mr.match_id = ?
            UNION ALL
            SELECT mr.team_a_player2_id AS player_id, p2.name AS player_name, mr.score_a AS score
            FROM match_rounds mr
            JOIN match_players p2 ON p2.id = mr.team_a_player2_id
            WHERE mr.match_id = ? AND mr.team_a_player2_id != 0
            UNION ALL
            SELECT mr.team_b_player1_id AS player_id, p1.name AS player_name, mr.score_b AS score
            FROM match_rounds mr
            JOIN match_players p1 ON p1.id = mr.team_b_player1_id
            WHERE mr.match_id = ?
            UNION ALL
            SELECT mr.team_b_player2_id AS player_id, p2.name AS player_name, mr.score_b AS score
            FROM match_rounds mr
            JOIN match_players p2 ON p2.id = mr.team_b_player2_id
            WHERE mr.match_id = ? AND mr.team_b_player2_id != 0
        ), player_totals AS (
            SELECT player_id, player_name, SUM(score) AS total_score
            FROM player_scores
            GROUP BY player_id, player_name
        )
        SELECT
            player_name AS name,
            total_score AS score,
            DENSE_RANK() OVER (ORDER BY total_score DESC, player_name ASC) AS global_rank
        FROM player_totals
        ORDER BY global_rank ASC, name ASC
    `

	err := r.db.Raw(query, matchID, matchID, matchID, matchID).Scan(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	return entries, len(entries), nil
}

func (r *leaderboardRepository) GetEventDetail(matchID string) (*models.EventDetail, error) {
	var detail models.EventDetail
	err := r.db.Table("match_events").
		Where("match_events.id = ?", matchID).
		First(&detail).Error

	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *leaderboardRepository) GetScoresByMatchID(matchID string) ([]models.ScoreEntry, error) {
	var entries []models.ScoreEntry

	query := `
        WITH individual_scores AS (
            SELECT mr.round_number, 'Team A' AS team_name, p1.name AS player_name, mr.score_a AS score
            FROM match_rounds mr
            JOIN match_players p1 ON p1.id = mr.team_a_player1_id
            WHERE mr.match_id = ?
            UNION ALL
            SELECT mr.round_number, 'Team A' AS team_name, p2.name AS player_name, mr.score_a AS score
            FROM match_rounds mr
            JOIN match_players p2 ON p2.id = mr.team_a_player2_id
            WHERE mr.match_id = ? AND mr.team_a_player2_id != 0
            UNION ALL
            SELECT mr.round_number, 'Team B' AS team_name, p1.name AS player_name, mr.score_b AS score
            FROM match_rounds mr
            JOIN match_players p1 ON p1.id = mr.team_b_player1_id
            WHERE mr.match_id = ?
            UNION ALL
            SELECT mr.round_number, 'Team B' AS team_name, p2.name AS player_name, mr.score_b AS score
            FROM match_rounds mr
            JOIN match_players p2 ON p2.id = mr.team_b_player2_id
            WHERE mr.match_id = ? AND mr.team_b_player2_id != 0
        )
        SELECT round_number, team_name, player_name, score
        FROM individual_scores
        ORDER BY round_number ASC, score DESC, player_name ASC
    `

	err := r.db.Raw(query, matchID, matchID, matchID, matchID).Scan(&entries).Error

	return entries, err
}

// GET LEADERBOARD BY SESSION
// func (r *leaderboardRepository) GetLeaderboardBySession(session string, page, pageSize int) ([]models.LeaderboardEntry, int, error) {
// 	var entries []models.LeaderboardEntry
// 	var total int64

// 	// Hitung total entries untuk pagination
// 	err := r.db.Table("match_players").
// 		Joins("JOIN users ON match_players.user_id = users.id").
// 		Joins("JOIN matches ON match_players.match_id = matches.id").
// 		Where("matches.session = ?", session).
// 		Count(&total).Error

// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Query untuk mendapatkan leaderboard berdasarkan temp_score di match_players dengan pagination
// 	err = r.db.Table("match_players").
// 		Select("users.id as user_id, users.username, SUM(match_players.temp_score) as score").
// 		Joins("JOIN users ON match_players.user_id = users.id").
// 		Joins("JOIN matches ON match_players.match_id = matches.id").
// 		Where("matches.session = ?", session).
// 		Group("users.id").
// 		Order("score DESC").
// 		Limit(pageSize).
// 		Offset((page - 1) * pageSize).
// 		Scan(&entries).Error

// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return entries, int(total), nil
// }
