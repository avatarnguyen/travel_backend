package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Place struct {
	ID               primitive.ObjectID `bson:"_id"`
	Name             *string            `json:"name"`
	Description      *string            `json:"description"`
	Street           *string            `json:"street"`
	City             *string            `json:"city"`
	Country          *string            `json:"country"`
	PostalCode       *string            `json:"postalCode"`
	AddressString    *string            `json:"addressString"`
	TripAdvisorID    *string            `json:"tripAdvisorId"`
	FourSquareID     *string            `json:"fourSquareId"`
	SearchQueries    []string           `json:"searchQueries"`
	SearchCategories []string           `json:"searchCategories"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
}
