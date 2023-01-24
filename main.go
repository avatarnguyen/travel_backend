package main

import (
	"log"
	"os"

	"github.com/avatarnguyen/travel_backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

// @title Travel Buddy
// @version 0.1
// @description Golang backend API using Fiber and MongoDB for Travel Buddy Mobile App.
// @contact.name Anh Nguyen
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.SearchCityRoutes(app)
	routes.IternaryRoutes(app)

	// attach swagger
	//config.AddSwaggerRoutes(app)

	err1 := app.Listen(":" + port)
	if err1 != nil {
		log.Panic("not able to start server")
	}
}
