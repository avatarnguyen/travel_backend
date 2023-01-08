package routes

import (
	"github.com/gofiber/fiber/v2"

	controller "github.com/avatarnguyen/travel_backend/controllers"
)

func AuthRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Post("users/signup", controller.SignUp())
	incomingRoutes.Post("users/login", controller.Login())
}
