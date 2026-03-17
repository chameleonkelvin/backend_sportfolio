package repositories

import (
    "scoring_app/models"
    "gorm.io/gorm"
)

type BannerEventRepository interface {
    GetAll() ([]models.BannerEvent, error)
    GetByID(id string) (*models.BannerEvent, error)
    GetByUserID(userID string) ([]models.BannerEvent, error)
    Create(event *models.BannerEvent) error
    Update(event *models.BannerEvent) error
    Delete(id string) error
}

type bannerEventRepository struct {
    db *gorm.DB
}

func NewBannerEventRepository(db *gorm.DB) BannerEventRepository {
    return &bannerEventRepository{db: db}
}

func (r *bannerEventRepository) GetAll() ([]models.BannerEvent, error) {
    var events []models.BannerEvent
    err := r.db.Find(&events).Error
    return events, err
}

func (r *bannerEventRepository) GetByID(id string) (*models.BannerEvent, error) {
    var event models.BannerEvent
    err := r.db.First(&event, "id = ?", id).Error
    return &event, err
}

func (r *bannerEventRepository) GetByUserID(userID string) ([]models.BannerEvent, error) {
    var events []models.BannerEvent
    err := r.db.Where("user_id = ?", userID).Find(&events).Error
    return events, err
}

func (r *bannerEventRepository) Create(event *models.BannerEvent) error {
    return r.db.Create(event).Error
}

func (r *bannerEventRepository) Update(event *models.BannerEvent) error {
    return r.db.Save(event).Error
}

func (r *bannerEventRepository) Delete(id string) error {
    return r.db.Delete(&models.BannerEvent{}, "id = ?", id).Error
}