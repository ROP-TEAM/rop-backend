package database

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&models.Company{},
		&models.User{},
		&models.DistributionPoint{},
		&models.Car{},
		&models.TagSkill{},
		&models.CarTagSkill{},
		&models.Order{},
		&models.OrderTagSkill{},
		&models.Client{},
		&models.Plan{},
		&models.Route{},
		&models.Route{})

	defaultSkills := []models.TagSkill{
		{Name: "อาหารสด"},
		{Name: "อาหารแช่แข็ง"},
		{Name: "Fragile (สินค้าแตกหักง่าย)"},
		{Name: "สินค้าขนาดใหญ่"},
		{Name: "สารเคมี/วัตถุอันตราย"},
	}

	for _, skill := range defaultSkills {
		db.Where(models.TagSkill{Name: skill.Name}).FirstOrCreate(&skill)
	}

	return db, nil
}
