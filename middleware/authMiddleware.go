package middleware

import (
	"github.com/avatarnguyen/travel_backend/helpers"

	"github.com/gofiber/fiber/v2"
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {

		clientToken := c.Get("token")
		if clientToken == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No Authorization header provided",
			})
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		return c.Next()
	}
}
