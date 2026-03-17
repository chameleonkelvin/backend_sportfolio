package services

import (
	"errors"
	"fmt"
	"scoring_app/models"
	"scoring_app/repositories"
	"time"

	"github.com/google/uuid"
)

type MatchEventService interface {
	Create(userID, accountTypeID, name string, totalCourts int, gameType, location string, playDate time.Time, totalPlayers int, teamType string) (*models.MatchEvent, error)
	GetByID(id string) (*models.MatchEvent, error)
	GetAll(page, pageSize int) ([]models.MatchEvent, int64, error)
	GetByUserID(userID string) ([]models.MatchEvent, error)
	Update(id, userID, name string, totalCourts int, gameType, location string, playDate time.Time, totalPlayers int, teamType string) (*models.MatchEvent, error)
	Delete(id, userID string) error
}

type matchEventService struct {
	repo repositories.MatchEventRepository
}

func NewMatchEventService(repo repositories.MatchEventRepository) MatchEventService {
	return &matchEventService{repo: repo}
}

func (s *matchEventService) Create(userID, accountTypeID, name string, totalCourts int, gameType, location string, playDate time.Time, totalPlayers int, teamType string) (*models.MatchEvent, error) {
	// Validate total players
	if totalPlayers < 4 {
		return nil, errors.New("total players must be at least 4")
	}

	// Validate total players must be multiple of 4
	if totalPlayers%4 != 0 {
		return nil, errors.New("total players must be a multiple of 4")
	}

	// Calculate expected players based on courts (2 courts = 8 players, 3 courts = 12 players, etc.)
	expectedPlayers := totalCourts * 4
	if totalPlayers != expectedPlayers {
		return nil, fmt.Errorf("for %d court(s), total players must be %d (got %d)", totalCourts, expectedPlayers, totalPlayers)
	}

	// Business rule: If courts > 2, only admin can create
	if totalCourts > 2 && accountTypeID != "admin" {
		return nil, errors.New("only admin can create match events with more than 2 courts")
	}

	// Create match event
	matchEvent := &models.MatchEvent{
		ID:           uuid.New().String(),
		UserID:       userID,
		Name:         name,
		TotalCourts:  totalCourts,
		GameType:     gameType,
		Location:     location,
		PlayDate:     playDate,
		TotalPlayers: totalPlayers,
		TeamType:     teamType,
	}

	err := s.repo.Create(matchEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to create match event: %w", err)
	}

	// Load with relationships
	matchEvent, _ = s.repo.FindByIDWithDetails(matchEvent.ID)

	return matchEvent, nil
}

func (s *matchEventService) GetByID(id string) (*models.MatchEvent, error) {
	matchEvent, err := s.repo.FindByIDWithDetails(id)
	if err != nil {
		return nil, fmt.Errorf("match event not found")
	}
	return matchEvent, nil
}

func (s *matchEventService) GetAll(page, pageSize int) ([]models.MatchEvent, int64, error) {
	matchEvents, total, err := s.repo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get match events: %w", err)
	}
	return matchEvents, total, nil
}

func (s *matchEventService) GetByUserID(userID string) ([]models.MatchEvent, error) {
	matchEvents, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match events: %w", err)
	}
	return matchEvents, nil
}

func (s *matchEventService) Update(id, userID, name string, totalCourts int, gameType, location string, playDate time.Time, totalPlayers int, teamType string) (*models.MatchEvent, error) {
	// Validate total players
	if totalPlayers < 4 {
		return nil, errors.New("total players must be at least 4")
	}

	// Validate total players must be multiple of 4
	if totalPlayers%4 != 0 {
		return nil, errors.New("total players must be a multiple of 4")
	}

	// Calculate expected players based on courts
	expectedPlayers := totalCourts * 4
	if totalPlayers != expectedPlayers {
		return nil, fmt.Errorf("for %d court(s), total players must be %d (got %d)", totalCourts, expectedPlayers, totalPlayers)
	}

	matchEvent, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("match event not found")
	}

	// Check ownership
	if matchEvent.UserID != userID {
		return nil, errors.New("you don't have permission to update this match event")
	}

	// Update fields
	matchEvent.Name = name
	matchEvent.TotalCourts = totalCourts
	matchEvent.GameType = gameType
	matchEvent.Location = location
	matchEvent.PlayDate = playDate
	matchEvent.TotalPlayers = totalPlayers
	matchEvent.TeamType = teamType

	err = s.repo.Update(matchEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to update match event: %w", err)
	}

	// Load with relationships
	matchEvent, _ = s.repo.FindByIDWithDetails(matchEvent.ID)

	return matchEvent, nil
}

func (s *matchEventService) Delete(id, userID string) error {
	matchEvent, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("match event not found")
	}

	// Check ownership
	if matchEvent.UserID != userID {
		return errors.New("you don't have permission to delete this match event")
	}

	err = s.repo.Delete(matchEvent.ID)
	if err != nil {
		return fmt.Errorf("failed to delete match event: %w", err)
	}

	return nil
}
