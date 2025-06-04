package repository

import (
	"own-paynet/models"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	Create(company *models.Company) error
	Update(company *models.Company) error
	FindByID(id uint) (*models.Company, error)
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db}
}

func (r *companyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) Update(company *models.Company) error {
	return r.db.Model(&models.Company{}).Where("id = ?", company.ID).Updates(company).Error
}
func (r *companyRepository) FindByID(id uint) (*models.Company, error) {
	var company models.Company
	result := r.db.First(&company, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &company, nil
}
