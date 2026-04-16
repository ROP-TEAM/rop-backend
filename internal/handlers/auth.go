package handlers

import (
	"ROP_Backend/internal/services"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// GoogleLogin godoc
// @Summary Login with Google
// @Description Verify Google ID Token and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param body body object{id_token=string} true "Google ID Token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/google [post]
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
	fmt.Println(user.Email)

	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}
