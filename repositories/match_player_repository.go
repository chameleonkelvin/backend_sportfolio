package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type MatchPlayerRepository interface {
	Create(player *models.MatchPlayer) error
	CreateBatch(players []models.MatchPlayer) error
	FindByID(id uint) (*models.MatchPlayer, error)
	FindByMatchID(matchID string) ([]models.MatchPlayer, error)
	Update(player *models.MatchPlayer) error
	Delete(id uint) error
	DeleteByMatchID(matchID string) error
}

type matchPlayerRepository struct {
	db *gorm.DB
}

func NewMatchPlayerRepository(db *gorm.DB) MatchPlayerRepository {
	return &matchPlayerRepository{db: db}
}

func (r *matchPlayerRepository) Create(player *models.MatchPlayer) error {
	return r.db.Create(player).Error
}

func (r *matchPlayerRepository) CreateBatch(players []models.MatchPlayer) error {
	return r.db.Create(&players).Error
}

func (r *matchPlayerRepository) FindByID(id uint) (*models.MatchPlayer, error) {
	var player models.MatchPlayer
	err := r.db.First(&player, id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *matchPlayerRepository) FindByMatchID(matchID string) ([]models.MatchPlayer, error) {
	var players []models.MatchPlayer
	err := r.db.Where("match_id = ?", matchID).
		Order("temp_score DESC, name ASC").
		Find(&players).Error
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (r *matchPlayerRepository) Update(player *models.MatchPlayer) error {
	return r.db.Save(player).Error
}

func (r *matchPlayerRepository) Delete(id uint) error {
	return r.db.Delete(&models.MatchPlayer{}, id).Error
}

func (r *matchPlayerRepository) DeleteByMatchID(matchID string) error {
	return r.db.Where("match_id = ?", matchID).Delete(&models.MatchPlayer{}).Error
}
