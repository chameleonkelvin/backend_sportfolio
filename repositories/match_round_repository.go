package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type MatchRoundRepository interface {
	Create(round *models.MatchRound) error
	GetAll(page int, pageSize int) ([]models.MatchRound, int64, error)
	FindByID(id uint) (*models.MatchRound, error)
	FindByIDAndCourt(id uint, roundCourt int) (*models.MatchRound, error)
	FindByIDWithPlayers(id uint) (*models.MatchRound, error)
	FindByMatchID(matchID string) ([]models.MatchRound, error)
	FindByMatchIDWithPlayers(matchID string) ([]models.MatchRound, error)
	FindByRoundNumber(matchID string, roundNumber int) ([]models.MatchRound, error)
	PatchScore(id uint, court int, addScoreA int, addScoreB int) error
	Update(round *models.MatchRound) error
	Delete(id uint) error
	DeleteByMatchID(matchID string) error
	CreateNextRound(lastRound *models.MatchRound) error
}

type matchRoundRepository struct {
	db *gorm.DB
}

func NewMatchRoundRepository(db *gorm.DB) MatchRoundRepository {
	return &matchRoundRepository{db: db}
}

// CREATE
func (r *matchRoundRepository) Create(round *models.MatchRound) error {
	// 1. Simpan ke database
	err := r.db.Create(round).Error
	if err != nil {
		return err
	}

	// 2. Ambil ulang data dengan relasi lengkap (Preload)
	// Menggunakan ID yang baru saja dibuat
	err = r.db.Preload("Match").
		Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		First(round, round.ID).Error

	return err
}

// GET ALL
func (r *matchRoundRepository) GetAll(page int, pageSize int) ([]models.MatchRound, int64, error) {

	var rounds []models.MatchRound
	var total int64

	offset := (page - 1) * pageSize

	// Hitung total data
	if err := r.db.Model(&models.MatchRound{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan limit + offset
	err := r.db.
		Preload("Match").
		Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&rounds).Error

	if err != nil {
		return nil, 0, err
	}

	return rounds, total, nil
}

// GET BY ID
func (r *matchRoundRepository) FindByID(id uint) (*models.MatchRound, error) {
	var rounds []models.MatchRound

	// Gunakan Preload untuk setiap field relasi yang didefinisikan di struct
	err := r.db.Preload("Match").
		Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		Order("match_id ASC, round_number ASC, court ASC").
		Find(&rounds).Error

	if err != nil {
		return nil, err
	}

	for _, round := range rounds {
		if round.ID == id {
			return &round, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

// GET BY ID AND COURT
func (r *matchRoundRepository) FindByIDAndCourt(id uint, court int) (*models.MatchRound, error) {
	var round models.MatchRound
	err := r.db.Where("id = ? AND court = ?", id, court).First(&round).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// GET BY ID WITH PLAYERS
func (r *matchRoundRepository) FindByIDWithPlayers(id uint) (*models.MatchRound, error) {
	var round models.MatchRound
	err := r.db.Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		First(&round, id).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

// GET BY MATCH ID
func (r *matchRoundRepository) FindByMatchID(matchID string) ([]models.MatchRound, error) {
	var rounds []models.MatchRound
	err := r.db.Where("match_id = ?", matchID).
		Order("round_number ASC, court ASC").
		Find(&rounds).Error
	if err != nil {
		return nil, err
	}
	return rounds, nil
}

// GET BY MATCH ID WITH PLAYERS
func (r *matchRoundRepository) FindByMatchIDWithPlayers(matchID string) ([]models.MatchRound, error) {
	var rounds []models.MatchRound
	err := r.db.Where("match_id = ?", matchID).
		Preload("Match").
		Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		Order("round_number ASC, court ASC").
		Find(&rounds).Error
	if err != nil {
		return nil, err
	}
	return rounds, nil
}

// GET BY ROUND NUMBER
func (r *matchRoundRepository) FindByRoundNumber(matchID string, roundNumber int) ([]models.MatchRound, error) {
	var rounds []models.MatchRound
	err := r.db.Where("match_id = ? AND round_number = ?", matchID, roundNumber).
		Preload("TeamAPlayer1").
		Preload("TeamAPlayer2").
		Preload("TeamBPlayer1").
		Preload("TeamBPlayer2").
		Order("court ASC").
		Find(&rounds).Error
	if err != nil {
		return nil, err
	}
	return rounds, nil
}

// UPDATE SCORES
func (r *matchRoundRepository) PatchScore(id uint, court int, addScoreA int, addScoreB int) error {

	updates := map[string]interface{}{}

	if addScoreA != 0 {
		updates["score_a"] = gorm.Expr("score_a + ?", addScoreA)
	}

	if addScoreB != 0 {
		updates["score_b"] = gorm.Expr("score_b + ?", addScoreB)
	}

	if len(updates) == 0 {
		return nil
	}

	result := r.db.Model(&models.MatchRound{}).
		Where("id = ? AND court = ?", id, court).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// INCREMENT ROUND NUMBER
func (r *matchRoundRepository) CreateNextRound(lastRound *models.MatchRound) error {
    // Buat struct baru berdasarkan data ronde terakhir
    newRound := models.MatchRound{
        MatchID:        lastRound.MatchID,
        Court:          lastRound.Court,
        RoundNumber:    lastRound.RoundNumber + 1, // Naikkan nomor ronde
        TeamAPlayer1ID: lastRound.TeamAPlayer1ID,
        TeamAPlayer2ID: lastRound.TeamAPlayer2ID,
        TeamBPlayer1ID: lastRound.TeamBPlayer1ID,
        TeamBPlayer2ID: lastRound.TeamBPlayer2ID,
        ScoreA:         0, // Reset skor
        ScoreB:         0, // Reset skor
        // Tambahkan field lain yang perlu di-copy (misal: TournamentID, dll)
    }

    return r.db.Create(&newRound).Error
}

// UPDATE
func (r *matchRoundRepository) Update(round *models.MatchRound) error {
	return r.db.Save(round).Error
}

// DELETE
func (r *matchRoundRepository) Delete(id uint) error {
	return r.db.Delete(&models.MatchRound{}, id).Error
}

// DELETE BY MATCH ID
func (r *matchRoundRepository) DeleteByMatchID(matchID string) error {
	return r.db.Where("match_id = ?", matchID).Delete(&models.MatchRound{}).Error
}
