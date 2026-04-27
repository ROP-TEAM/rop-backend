package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	return &user, err

}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByPhone(ctx context.Context, tel string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("tel = ?", tel).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
