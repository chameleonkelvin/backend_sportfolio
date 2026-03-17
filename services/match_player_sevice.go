package services

import (
	"errors"
	"strconv"

	"scoring_app/models"
	"scoring_app/repositories"
)

type MatchPlayerService interface {
	Create(player *models.MatchPlayer) error
	CreateBatch(players []models.MatchPlayer) error
	GetAll(page int, pageSize int) ([]models.MatchPlayer, int64, error)
	GetByMatchID(matchID string) ([]models.MatchPlayer, error)
	FindByID(id string) (*models.MatchPlayer, error) // 🔥 string
	FindByMatchID(matchID string) ([]models.MatchPlayer, error)
	Update(player *models.MatchPlayer) error
	Delete(id string) error // 🔥 string
	DeleteByMatchID(matchID string) error
}

type matchPlayerService struct {
	repo           repositories.MatchPlayerRepository
	matchEventRepo repositories.MatchEventRepository
}

func NewMatchPlayerService(
	repo repositories.MatchPlayerRepository,
	matchEventRepo repositories.MatchEventRepository,
) MatchPlayerService {
	return &matchPlayerService{
		repo:           repo,
		matchEventRepo: matchEventRepo,
	}
}

// CREATE
func (s *matchPlayerService) Create(player *models.MatchPlayer) error {

	_, err := s.matchEventRepo.FindByID(player.MatchID)
	if err != nil {
		return errors.New("match_event not found")
	}

	return s.repo.Create(player)
}

// BULK CREATE
func (s *matchPlayerService) CreateBatch(players []models.MatchPlayer) error {

	for _, player := range players {
		_, err := s.matchEventRepo.FindByID(player.MatchID)
		if err != nil {
			return errors.New("match_event not found")
		}
	}

	return s.repo.CreateBatch(players)
}

// GET ALL
func (s *matchPlayerService) GetAll(page int, pageSize int) ([]models.MatchPlayer, int64, error) {
	return s.repo.GetAll(page, pageSize)
}

// GET ALL BY MATCH ID
func (s *matchPlayerService) GetByMatchID(matchID string) ([]models.MatchPlayer, error) {
	return s.repo.FindByMatchID(matchID)
}

// FIND BY ID
func (s *matchPlayerService) FindByID(id string) (*models.MatchPlayer, error) {
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	return s.repo.FindByID(uint(idUint64))
}

// FIND BY MATCH ID
func (s *matchPlayerService) FindByMatchID(matchID string) ([]models.MatchPlayer, error) {
	return s.repo.FindByMatchID(matchID)
}

// UPDATE
func (s *matchPlayerService) Update(player *models.MatchPlayer) error {
	return s.repo.Update(player)
}

// DELETE BY ID
func (s *matchPlayerService) Delete(id string) error {
	idUint64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errors.New("invalid id format")
	}

	return s.repo.Delete(uint(idUint64))
}

// DELETE BY MATCH ID
func (s *matchPlayerService) DeleteByMatchID(matchID string) error {
	return s.repo.DeleteByMatchID(matchID)
}
