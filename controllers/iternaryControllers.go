package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/avatarnguyen/travel_backend/database"
	"github.com/avatarnguyen/travel_backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var iternaryCollection = database.OpenConnection(database.Client, "iternary")

func CreateIternary() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		iternary := new(models.Iternary)
		if err := c.BodyParser(iternary); err != nil {
			cancel()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		fmt.Println("Destinations: \n", iternary.Destinations)

		defer cancel()

		iternary.ID = primitive.NewObjectID()
		iternary.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		iternary.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		iternary.IternaryID = iternary.ID.Hex()

		for i := range iternary.Destinations {
			destination := iternary.Destinations[i]

			destination.ID = primitive.NewObjectID()
			destination.Place.ID = primitive.NewObjectID()
			destination.Place.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			destination.Place.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		}

		resultInsertedNumber, insertErr := iternaryCollection.InsertOne(ctx, iternary)

		if insertErr != nil {
			msg := fmt.Sprintln("Iternary could not be inserted. Error: " + insertErr.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": msg,
			})
		}

		return c.Status(fiber.StatusOK).JSON(resultInsertedNumber)
	}
}
