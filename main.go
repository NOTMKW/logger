package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/notmkw/log/internal/config"
	"github.com/notmkw/log/internal/routes/routes"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	routes.Setup(app)

	log.Printf("Starting %s on port %s...\n", cfg.AppName, cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))

}
