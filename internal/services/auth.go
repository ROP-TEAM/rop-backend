package services

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/models"
	"context"
	"errors"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

func (s *AuthService) GoogleLogin(idToken string) (string, models.User, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, s.cfg.GOOGLE_CLIENT_ID)
	// payload, err := idtoken.Validate(context.Background(), idToken, "") --test
	if err != nil {
		return "", models.User{}, errors.New("invalid google token")
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	googleID := payload.Subject

	if verified, ok := payload.Claims["email_verified"].(bool); !ok || !verified {
		return "", models.User{}, errors.New("email not verified")
	}

	var user models.User
	result := s.db.Where("google_id = ?", googleID).First(&user)

	if result.Error != nil {
		user = models.User{
			Name:     name,
			Email:    email,
			Pictures: picture,
			GoogleID: googleID,
		}
		s.db.Create(&user)
	}

	token, err := s.makeToken(user.ID)
	if err != nil {
		return "", models.User{}, err
	}

	return token, user, nil
}
