package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type MatchEventRepository interface {
	Create(matchEvent *models.MatchEvent) error
	FindByID(id string) (*models.MatchEvent, error)
	FindByIDWithDetails(id string) (*models.MatchEvent, error)
	FindAll(page, pageSize int) ([]models.MatchEvent, int64, error)
	FindByUserID(userID string) ([]models.MatchEvent, error)
	FindUpcoming() ([]models.MatchEvent, error)
	Update(matchEvent *models.MatchEvent) error
	Delete(id string) error
}

type matchEventRepository struct {
	db *gorm.DB
}

func NewMatchEventRepository(db *gorm.DB) MatchEventRepository {
	return &matchEventRepository{db: db}
}

func (r *matchEventRepository) Create(matchEvent *models.MatchEvent) error {
	return r.db.Create(matchEvent).Error
}

func (r *matchEventRepository) FindByID(id string) (*models.MatchEvent, error) {
	var matchEvent models.MatchEvent
	err := r.db.First(&matchEvent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &matchEvent, nil
}

func (r *matchEventRepository) FindByIDWithDetails(id string) (*models.MatchEvent, error) {
	var matchEvent models.MatchEvent
	err := r.db.Preload("User").
		Preload("Players").
		Preload("Rounds").
		First(&matchEvent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &matchEvent, nil
}

func (r *matchEventRepository) FindAll(page, pageSize int) ([]models.MatchEvent, int64, error) {
	var matchEvents []models.MatchEvent
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&models.MatchEvent{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").
		Offset(offset).
		Limit(pageSize).
		Order("play_date DESC").
		Find(&matchEvents).Error
	if err != nil {
		return nil, 0, err
	}

	return matchEvents, total, nil
}

func (r *matchEventRepository) FindByUserID(userID string) ([]models.MatchEvent, error) {
	var matchEvents []models.MatchEvent
	err := r.db.Where("user_id = ?", userID).
		Order("play_date DESC").
		Find(&matchEvents).Error
	if err != nil {
		return nil, err
	}
	return matchEvents, nil
}

func (r *matchEventRepository) FindUpcoming() ([]models.MatchEvent, error) {
	var matchEvents []models.MatchEvent
	err := r.db.Where("play_date > NOW()").
		Order("play_date ASC").
		Preload("User").
		Find(&matchEvents).Error
	if err != nil {
		return nil, err
	}
	return matchEvents, nil
}

func (r *matchEventRepository) Update(matchEvent *models.MatchEvent) error {
	return r.db.Save(matchEvent).Error
}

func (r *matchEventRepository) Delete(id string) error {
	return r.db.Delete(&models.MatchEvent{}, "id = ?", id).Error
}
