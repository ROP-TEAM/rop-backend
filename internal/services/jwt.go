package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) makeToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT_SECRET))
}
