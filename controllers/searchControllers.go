package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/avatarnguyen/travel_backend/database"
	"github.com/avatarnguyen/travel_backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var placeCollection = database.OpenConnection(database.Client, "places")

func Search(query string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		key := os.Getenv("TRIPADVISOR_KEY")

		searchQuery := "lisbon"
		searchCategory := "attractions"

		url := "https://api.content.tripadvisor.com/api/v1/location/search?key=" +
			key +
			"&searchQuery=" + searchQuery +
			"&category=" + searchCategory +
			"&language=en"

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			cancel()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			cancel()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		defer cancel()

		var jsonRes map[string]interface{}
		_ = json.Unmarshal(body, &jsonRes)

		arr := jsonRes["data"].([]interface{})

		for i := range arr {
			place := new(models.Place)

			place.ID = primitive.NewObjectID()

			place.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			place.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			queries := []string{searchQuery}
			place.SearchQueries = queries

			categories := []string{searchCategory}
			place.SearchCategories = categories

			arr_item := arr[i].(map[string]interface{})

			// Place
			name := arr_item["name"].(string)
			place.Name = &name

			trip_advisor_id := arr_item["location_id"].(string)
			place.TripAdvisorID = &trip_advisor_id

			// Address
			address_obj := arr_item["address_obj"].(map[string]interface{})

			address_string := address_obj["address_string"].(string)
			place.AddressString = &address_string

			_, insertErr := placeCollection.InsertOne(ctx, place)
			if insertErr != nil {
				cancel()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": insertErr.Error(),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(len(arr))
	}
}
