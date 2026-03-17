package services

import (
	"scoring_app/models"
	"scoring_app/repositories"
	"time"

	"github.com/google/uuid"
)

// Gunakan struct yang sama atau bedakan sesuai kebutuhan controller
type CreateBannerEventInput struct {
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// Tambahkan ini agar tidak error "undefined: services.UpdateBannerEventInput"
type UpdateBannerEventInput struct {
	Title       string    `json:"title"`
	Location    string    `json:"location"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	IsActive    *bool     `json:"is_active"` // Gunakan pointer agar bisa handle nilai false
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type BannerEventService interface {
	GetAllBannerEvents() ([]models.BannerEvent, error)
	GetBannerEventByID(id string) (*models.BannerEvent, error)
	GetBannerEventByUserId(userID string) ([]models.BannerEvent, error)
	CreateBannerEvent(input CreateBannerEventInput) (*models.BannerEvent, error)
	UpdateBannerEvent(id string, input UpdateBannerEventInput) (*models.BannerEvent, error)
	DeleteBannerEvent(id string) error
}

type bannerEventService struct {
	repo repositories.BannerEventRepository
}

func NewBannerEventService(repo repositories.BannerEventRepository) BannerEventService {
	return &bannerEventService{repo: repo}
}

func (s *bannerEventService) CreateBannerEvent(input CreateBannerEventInput) (*models.BannerEvent, error) {
	event := &models.BannerEvent{
		ID:          uuid.New().String(),
		UserID:      input.UserID,
		Title:       input.Title,
		Location:    input.Location,
		Image:       input.Image,
		Description: input.Description,
		IsActive:    input.IsActive,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}
	err := s.repo.Create(event)
	return event, err
}

// Update fungsi agar menerima UpdateBannerEventInput
func (s *bannerEventService) UpdateBannerEvent(id string, input UpdateBannerEventInput) (*models.BannerEvent, error) {
	event, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update field jika dikirim (tidak kosong)
	if input.Title != "" {
		event.Title = input.Title
	}
	if input.Location != "" {
		event.Location = input.Location
	}
	if input.Image != "" {
		event.Image = input.Image
	}
	if input.IsActive != nil {
		event.IsActive = *input.IsActive
	}
	if !input.StartDate.IsZero() {
		event.StartDate = input.StartDate
	}
	if !input.EndDate.IsZero() {
		event.EndDate = input.EndDate
	}

	err = s.repo.Update(event)
	return event, err
}

func (s *bannerEventService) GetAllBannerEvents() ([]models.BannerEvent, error) {
	return s.repo.GetAll()
}

func (s *bannerEventService) GetBannerEventByID(id string) (*models.BannerEvent, error) {
	return s.repo.GetByID(id)
}

func (s *bannerEventService) GetBannerEventByUserId(userID string) ([]models.BannerEvent, error) {
	return s.repo.GetByUserID(userID)
}

func (s *bannerEventService) DeleteBannerEvent(id string) error {
	return s.repo.Delete(id)
}
