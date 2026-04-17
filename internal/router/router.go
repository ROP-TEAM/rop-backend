package router

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/handlers"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	_ "ROP_Backend/docs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	app.Get("/swagger/*", adaptor.HTTPHandler(
		httpSwagger.WrapHandler,
	))

	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	api := app.Group("/api")

	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Get("/test", middleware.Protected(cfg), handlers.Test)

	return app
}
