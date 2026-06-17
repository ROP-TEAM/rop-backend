package router

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/handlers"
	"ROP_Backend/internal/middleware"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/services"

	_ "ROP_Backend/docs"

	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(middleware.Cors())

	app.Use(middleware.RateLimiter())

	// app.Get("/api/swagger/*", adaptor.HTTPHandler(
	// 	httpSwagger.WrapHandler,
	// ))

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

	matrixService, err := services.NewMatrixService(cfg)
	if err != nil {
		panic(err)
	}
	matrixHandler := handlers.NewMatrixHandler(matrixService)

	slv, solverCleanup, err := buildSolver(cfg.SOLVER_BINARY_PATH)
	if err != nil {
		panic(fmt.Sprintf("failed to start solver: %v", err))
	}
	if solverCleanup != nil {
		app.Hooks().OnPreShutdown(func() error {
			solverCleanup()
			return nil
		})
	}

	planRepo := repository.NewPlanRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	tagRepo := repository.NewTagSkillRepository(db)

	planningService := services.NewPlanningService(
		matrixService,
		slv,
		planRepo,
		vehicleRepo,
		orderRepo,
		tagRepo,
	)

	planningHandler := handlers.NewPlanningHandler(planningService)

	wastePlanningService := services.NewWastePlanningService(slv)
	wasteHandler := handlers.NewWasteHandler(wastePlanningService)

	api := app.Group("/api")

	api.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler(
		httpSwagger.URL("/api/swagger/doc.json"),
	)))

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

	api.Post("/plans", middleware.Protected(cfg), planHandler.Create)
	api.Patch("/plans/name/:id", middleware.Protected(cfg), planHandler.UpdateNameByID)
	api.Delete("/plans/:id", middleware.Protected(cfg), planHandler.DeleteByID)
	api.Post("/plans/:id", middleware.Protected(cfg), planHandler.DuplicateByID)
	api.Get("/plans/:id", middleware.Protected(cfg), planHandler.GetByID)
	api.Get("/plans/", middleware.Protected(cfg), planHandler.GetPlans)

	api.Post("/matrix", middleware.Protected(cfg), matrixHandler.BuildMatrix)

	api.Post("/optimize", planningHandler.Optimize)

	api.Post("/waste/plan", wasteHandler.Plan)

	//test route

	api.Get("/test", middleware.Protected(cfg), handlers.Test)

	api.Post("/auth/otp", middleware.OTPLimiter(), handlers.TestOTP)

	return app
}
