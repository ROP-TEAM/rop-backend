package handlers

import (
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) GoogleLogin(c fiber.Ctx) error {
	var body struct {
		IDToken string `json:"id_token"`
	}

	if err := c.Bind().Body(&body); err != nil || body.IDToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing id_token"})
	}

	token, user, err := h.service.GoogleLogin(body.IDToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}
