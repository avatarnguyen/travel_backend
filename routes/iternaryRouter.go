package routes

import (
	controller "github.com/avatarnguyen/travel_backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func IternaryRoutes(incomingRoutes *fiber.App) {
	incomingRoutes.Group("iternary")
	incomingRoutes.Post("create", controller.CreateIternary())
	incomingRoutes.Get(":id", controller.GetIternary())
	incomingRoutes.Patch(":id/edit", controller.EditIterary())
	incomingRoutes.Patch(":id/delete", controller.DeleteIternary())
}
