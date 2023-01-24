package routes

import (
	controller "github.com/avatarnguyen/travel_backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SearchCityRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Get("search/city", controller.Search("test"))

}
