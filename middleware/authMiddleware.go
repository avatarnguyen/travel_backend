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
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("uid", claims.UID)
		c.Set("user_type", claims.UserType)
		return c.Next()
	}
}
