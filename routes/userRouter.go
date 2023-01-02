package routes

import (
	controller "github.com/avatarnguyen/travel_backend/controllers"
	middleware "github.com/avatarnguyen/travel_backend/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoute *gin.Engine) {
	incomingRoute.Use(middleware.Authenticate())
	incomingRoute.GET("/users", controller.GetUsers())
	incomingRoute.GET("/users/:user_id", controller.GetUser())
}
