package repositories

import (
	"scoring_app/models"

	"gorm.io/gorm"
)

type AccountTypeRepository interface {
	Create(accountType *models.AccountType) error
	FindByID(id string) (*models.AccountType, error)
	FindAll() ([]models.AccountType, error)
	Update(accountType *models.AccountType) error
	Delete(id string) error
}

type accountTypeRepository struct {
	db *gorm.DB
}

func NewAccountTypeRepository(db *gorm.DB) AccountTypeRepository {
	return &accountTypeRepository{db: db}
}

func (r *accountTypeRepository) Create(accountType *models.AccountType) error {
	return r.db.Create(accountType).Error
}

func (r *accountTypeRepository) FindByID(id string) (*models.AccountType, error) {
	var accountType models.AccountType
	err := r.db.First(&accountType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &accountType, nil
}

func (r *accountTypeRepository) FindAll() ([]models.AccountType, error) {
	var accountTypes []models.AccountType
	err := r.db.Order("name ASC").Find(&accountTypes).Error
	if err != nil {
		return nil, err
	}
	return accountTypes, nil
}

func (r *accountTypeRepository) Update(accountType *models.AccountType) error {
	return r.db.Save(accountType).Error
}

func (r *accountTypeRepository) Delete(id string) error {
	return r.db.Delete(&models.AccountType{}, "id = ?", id).Error
}
