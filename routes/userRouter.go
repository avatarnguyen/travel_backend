package routes

import (
	controller "github.com/avatarnguyen/travel_backend/controllers"
	middleware "github.com/avatarnguyen/travel_backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoute *fiber.App) {
	incomingRoute.Use(middleware.Authenticate())
	incomingRoute.Get("/users", controller.GetUsers())
	incomingRoute.Get("/users/:user_id", controller.GetUser())
}
