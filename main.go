package main

import (
	_ "ROP_Backend/docs"
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/database"
	"ROP_Backend/internal/router"
	"log"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("cannot connect:", err)
	}

	app := router.Setup(db, cfg)
	log.Fatal(app.Listen(":" + cfg.APP_PORT))
}
