package services

import (
	"ROP_Backend/internal/repository"

	"gorm.io/gorm"
)

type CompanyService struct {
	companyRepository *repository.CompanyRepository
}

func NewCompanyService(db *gorm.DB) *CompanyService {
	companyRepo := repository.NewCompanyRepository(db)
	return &CompanyService{
		companyRepository: companyRepo,
	}
}
