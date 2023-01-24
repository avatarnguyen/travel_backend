package routes

import (
	controller "github.com/avatarnguyen/travel_backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func IternaryRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("iternary/create", controller.CreateIternary())
}
