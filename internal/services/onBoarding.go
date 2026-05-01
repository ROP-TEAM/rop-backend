package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

func (s *UserService) Onboarding(userID uint, payload dto.OnboardingPayload) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		var alley *string
		if payload.Alley != nil {
			alley = payload.Alley
		}

		company := models.Company{
			Name:        payload.CompanyName,
			Type:        payload.CompanyType,
			Province:    payload.Province,
			District:    payload.District,
			SubDistrict: payload.SubDistrict,
			Address:     payload.Address,
			Alley:       alley,
			PostalCode:  payload.PostalCode,
		}

		if err := tx.Create(&company).Error; err != nil {
			return err
		}

		user.CompanyID = &company.ID
		user.IsValidated = false

		if err := tx.Save(user).Error; err != nil {
			return err
		}

		return nil
	})
}
