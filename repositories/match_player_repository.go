package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type MatchPlayerRepository interface {
	GetAll(page int, pageSize int) ([]models.MatchPlayer, int64, error)
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

func (r *matchPlayerRepository) GetAll(page int, pageSize int) ([]models.MatchPlayer, int64, error) {
	var players []models.MatchPlayer
	var total int64

	offset := (page - 1) * pageSize

	// Hitung total data
	if err := r.db.Model(&models.MatchPlayer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan limit + offset
	err := r.db.
		Limit(pageSize).
		Offset(offset).
		Order("temp_score DESC, name ASC").
		Find(&players).Error

	if err != nil {
		return nil, 0, err
	}

	return players, total, nil
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
