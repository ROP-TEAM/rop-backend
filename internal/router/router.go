package router

import (
	"log"

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

	"github.com/ROP-TEAM/rop-algorithm/gmap"
	"github.com/ROP-TEAM/rop-algorithm/solver"
	grpcsolver "github.com/ROP-TEAM/rop-algorithm/solver/grpc"
	solverprocess "github.com/ROP-TEAM/rop-algorithm/solver/process"
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

	matrix, err := gmap.NewGoogleMapsMatrix(
		cfg.GOOGLE_MAPS_API_KEY,
		gmap.WithInMemoryCache(gmap.DevMatrixCacheConfig()),
	)
	if err != nil {
		log.Fatal("init matrix:", err)
	}

	var s solver.Solver = solver.NewStub()
	if cfg.SOLVER_BINARY_PATH != "" {
		handle, err := solverprocess.Start(cfg.SOLVER_BINARY_PATH)
		if err != nil {
			log.Fatal("start solver:", err)
		}
		s = grpcsolver.New(handle.Conn)
	}

	planningService := services.NewPlanningService(db, matrix, s)
	_ = planningService

	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	api := app.Group("/api")

	api.Post("/auth/google", authHandler.GoogleLogin)

	api.Get("/test", middleware.Protected(cfg), handlers.Test)

	api.Post("/onboarding",
		middleware.Protected(cfg),
		userHandler.Onboarding,
	)

	return app
}
