package services

import (
	"ROP_Backend/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
	db          *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo:    repository.NewUserRepository(db),
		companyRepo: repository.NewCompanyRepository(db),
		db:          db,
	}
}
