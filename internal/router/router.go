package router

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/handlers"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New()

	authService := services.NewAuthService(db, cfg)

	authHandler := handlers.NewAuthHandler(authService)

	api := app.Group("/api")

	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Get("/me", middleware.Protected(cfg), func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		return c.JSON(fiber.Map{"user_id": userID})
	})

	return app
}
