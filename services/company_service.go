package services

import (
	"fmt"
	"own-paynet/models"
	"own-paynet/repository"
)

type CompanyService interface {
	CreateCompany(company *models.Company) error
	UpdateCompany(id uint, updated *models.Company) (*models.Company, error)
}

type companyService struct {
	repo repository.CompanyRepository
}

func NewCompanyService(r repository.CompanyRepository) CompanyService {
	return &companyService{r}
}

func (s *companyService) CreateCompany(company *models.Company) error {
	return s.repo.Create(company)
}

func (s *companyService) UpdateCompany(id uint, updated *models.Company) (*models.Company, error) {
	company, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("record not found")
	}

	if updated.CompanyName != "" {
		company.CompanyName = updated.CompanyName
	}
	if updated.CompanyDetails != "" {
		company.CompanyDetails = updated.CompanyDetails
	}

	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	// Optional: reload from DB to get updated timestamps etc.
	return company, nil
}
