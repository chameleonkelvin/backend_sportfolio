package services

import (
	"errors"
	"fmt"
	"scoring_app/models"
	"scoring_app/repositories"

	"github.com/google/uuid"
)

type AccountTypeService interface {
	Create(name, description string) (*models.AccountType, error)
	GetByID(id string) (*models.AccountType, error)
	GetAll() ([]models.AccountType, error)
	Update(id, name, description string) (*models.AccountType, error)
	Delete(id string) error
}

type accountTypeService struct {
	repo repositories.AccountTypeRepository
}

func NewAccountTypeService(repo repositories.AccountTypeRepository) AccountTypeService {
	return &accountTypeService{repo: repo}
}

func (s *accountTypeService) Create(name, description string) (*models.AccountType, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	accountType := &models.AccountType{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}

	err := s.repo.Create(accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to create account type: %w", err)
	}

	return accountType, nil
}

func (s *accountTypeService) GetByID(id string) (*models.AccountType, error) {
	accountType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("account type not found")
	}
	return accountType, nil
}

func (s *accountTypeService) GetAll() ([]models.AccountType, error) {
	accountTypes, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get account types: %w", err)
	}
	return accountTypes, nil
}

func (s *accountTypeService) Update(id, name, description string) (*models.AccountType, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	accountType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("account type not found")
	}

	accountType.Name = name
	accountType.Description = description

	err = s.repo.Update(accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to update account type: %w", err)
	}

	return accountType, nil
}

func (s *accountTypeService) Delete(id string) error {
	accountType, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("account type not found")
	}

	err = s.repo.Delete(accountType.ID)
	if err != nil {
		return fmt.Errorf("failed to delete account type: %w", err)
	}

	return nil
}
