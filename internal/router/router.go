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

	vehicleService := services.NewVehicleService(db)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	orderService := services.NewOrderService(db)
	orderHandler := handlers.NewOrderHandler(orderService)

	tagSkillService := services.NewTagSkillService(db)
	tagSkillHandler := handlers.NewTagSkillHandler(tagSkillService)

	planService := services.NewPlanService(db, cfg)
	planHandler := handlers.NewPlanHandler(planService)

	api := app.Group("/api")
	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Post("/onboarding", middleware.Protected(cfg), userHandler.Onboarding)

	api.Post("/onboarding/otp", middleware.Protected(cfg), middleware.OTPLimiter(), otpHandler.RequestOTP)
	api.Post("/onboarding/otp/verify", middleware.Protected(cfg), otpHandler.VerifyOTP) //verifyมันฟรีไม่ต้องใส่limiter ??

	api.Post("/vehicles", middleware.Protected(cfg), vehicleHandler.Create)
	api.Patch("/vehicles/:id", middleware.Protected(cfg), vehicleHandler.Patch)
	api.Delete("/vehicles", middleware.Protected(cfg), vehicleHandler.Delete)

	api.Post("/orders", middleware.Protected(cfg), orderHandler.Create)
	api.Patch("/orders/:id", middleware.Protected(cfg), orderHandler.Patch)
	api.Delete("/orders", middleware.Protected(cfg), orderHandler.Delete)

	api.Post("/skills", middleware.Protected(cfg), tagSkillHandler.Create)

	api.Post("/plan", middleware.Protected(cfg), planHandler.Create)
	api.Patch("/plan/name/:id", middleware.Protected(cfg), planHandler.UpdateNameByID)
	api.Delete("/plan/:id", middleware.Protected(cfg), planHandler.DeleteByID)
	api.Post("/plan/:id", middleware.Protected(cfg), planHandler.DuplicateByID)

	//test route

	api.Get("/test", middleware.Protected(cfg), handlers.Test)

	api.Post("/auth/otp", middleware.OTPLimiter(), handlers.TestOTP)

	return app
}
