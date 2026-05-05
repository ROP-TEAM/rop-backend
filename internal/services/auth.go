package services

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"context"
	"errors"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepository *repository.UserRepository
	cfg            *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	userRepository := repository.NewUserRepository(db)
	return &AuthService{
		userRepository: userRepository,
		cfg:            cfg,
	}
}

func (s *AuthService) GoogleLogin(idToken string) (string, *models.User, bool, error) {
	// payload, err := idtoken.Validate(context.Background(), idToken, s.cfg.GOOGLE_CLIENT_ID)
	payload, err := idtoken.Validate(context.Background(), idToken, "")
	if err != nil {
		return "", nil, false, errors.New("invalid google token")
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	googleID := payload.Subject

	if verified, ok := payload.Claims["email_verified"].(bool); !ok || !verified {
		return "", nil, false, errors.New("email not verified")
	}

	user, err := s.userRepository.FindByGoogleID(googleID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := models.User{
			Name:        name,
			Email:       email,
			GoogleID:    googleID,
			IsValidated: false,
		}

		err = s.userRepository.Create(&newUser)
		if err != nil {
			return "", nil, false, err
		}

		user = &newUser
	} else if err != nil {
		return "", nil, false, err
	}

	needOnboarding := user.CompanyID == nil

	token, err := s.makeToken(user.ID, user.Email)
	if err != nil {
		return "", nil, false, err
	}

	return token, user, needOnboarding, nil
}
