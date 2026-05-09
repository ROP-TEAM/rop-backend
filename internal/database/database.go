package database

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Company{},
		&models.Plan{},
		&models.User{},
		&models.TagSkill{},
		&models.Vehicle{},
		&models.VehicleTagSkill{},
		&models.Order{},
		&models.OrderTagSkill{},
		&models.Route{},
		&models.Otp{},
		models.MockOTP{},
	); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	return db, nil
}
