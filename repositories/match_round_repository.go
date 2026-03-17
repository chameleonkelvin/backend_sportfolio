package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type MatchRoundRepository interface {
	Create(round *models.MatchRound) error
	CreateBatch(rounds []models.MatchRound) error
	FindByID(id uint) (*models.MatchRound, error)
	FindByIDWithPlayers(id uint) (*models.MatchRound, error)
	FindByMatchID(matchID string) ([]models.MatchRound, error)
	FindByMatchIDWithPlayers(matchID string) ([]models.MatchRound, error)
	FindByRoundNumber(matchID string, roundNumber int) ([]models.MatchRound, error)
	Update(round *models.MatchRound) error
	Delete(id uint) error
	DeleteByMatchID(matchID string) error
}

type matchRoundRepository struct {
	db *gorm.DB
}

func NewMatchRoundRepository(db *gorm.DB) MatchRoundRepository {
	return &matchRoundRepository{db: db}
}

func (r *matchRoundRepository) Create(round *models.MatchRound) error {
	return r.db.Create(round).Error
}

func (r *matchRoundRepository) CreateBatch(rounds []models.MatchRound) error {
	return r.db.Create(&rounds).Error
}

func (r *matchRoundRepository) FindByID(id uint) (*models.MatchRound, error) {
	var round models.MatchRound
	err := r.db.First(&round, id).Error
	if err != nil {
		return nil, err
	}
	return &round, nil
}

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

func (r *matchRoundRepository) FindByMatchIDWithPlayers(matchID string) ([]models.MatchRound, error) {
	var rounds []models.MatchRound
	err := r.db.Where("match_id = ?", matchID).
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

func (r *matchRoundRepository) Update(round *models.MatchRound) error {
	return r.db.Save(round).Error
}

func (r *matchRoundRepository) Delete(id uint) error {
	return r.db.Delete(&models.MatchRound{}, id).Error
}

func (r *matchRoundRepository) DeleteByMatchID(matchID string) error {
	return r.db.Where("match_id = ?", matchID).Delete(&models.MatchRound{}).Error
}
