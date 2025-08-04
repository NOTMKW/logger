package main

import (
	"log"
)

func main () {
	cfg := config.Load()

	app := fiber.New()

	routes.Setup(app)

	log.Printf("Starting %s on port %s...\n", cfg.AppName, cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
	
}