package router

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/handlers"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/services"

	_ "ROP_Backend/docs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(middleware.Cors())

	app.Use(middleware.RateLimiter())

	app.Get("/swagger/*", adaptor.HTTPHandler(
		httpSwagger.WrapHandler,
	))

	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	otpService := services.NewOTPService(db, cfg)
	otpHandler := handlers.NewOTPHandler(otpService)

	api := app.Group("/api")
	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Post("/onboarding", middleware.Protected(cfg), userHandler.Onboarding)

	api.Post("/onboarding/otp", middleware.Protected(cfg), otpHandler.RequestOTP)
	api.Post("/onboarding/otp/verify", middleware.Protected(cfg), otpHandler.VerifyOTP)

	//test route

	api.Get("/test", middleware.Protected(cfg), handlers.Test)

	api.Post("/auth/otp", middleware.OTPLimiter(), handlers.TestOTP)

	api.Post("/onboarding",
		middleware.Protected(cfg),
		userHandler.Onboarding,
	)

	return app
}
