package services

import (
	"time"

	"ROP_Backend/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) makeToken(userID uint, email string) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT_SECRET))
}
