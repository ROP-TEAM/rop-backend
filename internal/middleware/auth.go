package middleware

import (
	"ROP_Backend/internal/config"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func Protected(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Missing Token"})
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid Token Format"})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT_SECRET), nil
		})

		if err != nil {
			switch err.Error() {
			case jwt.ErrTokenExpired.Error():
				return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
			case jwt.ErrTokenMalformed.Error():
				return c.Status(401).JSON(fiber.Map{"error": "Token malformed"})
			case jwt.ErrSignatureInvalid.Error():
				return c.Status(401).JSON(fiber.Map{"error": "Invalid token signature"})
			default:
				return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
			}
		}

		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token is not valid"})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func GetUser(c fiber.Ctx) *Claims {
	if claims, ok := c.Locals("claims").(*Claims); ok {
		return claims
	}
	return nil
}
