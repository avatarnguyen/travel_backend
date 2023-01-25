package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Iternary struct {
	ID           primitive.ObjectID `bson:"_id"`
	IternaryID   string             `json:"iternaryId"`
	Name         string             `json:"name" validate:"required"`
	Description  *string            `json:"description"`
	Note         *string            `json:"note"`
	City         *string            `json:"city"`
	Country      string             `json:"country" validate:"required"`
	CreatorID    string             `json:"creatorId" validate:"required"`
	OwnerID      string             `json:"ownerId" validate:"required"`
	CreatorName  *string            `json:"creator"`
	OwnerName    *string            `json:"owner"`
	Tags         []string           `json:"tags"`
	IsTemplate   bool               `json:"isTemplate"`
	Rating       int                `json:"rating"`
	Destinations []Destination      `json:"destinations"`
	StartDate    time.Time          `json:"startDate"`
	EndDate      time.Time          `json:"endDate"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}

type Destination struct {
	ID            primitive.ObjectID `bson:"_id"`
	City          *string            `json:"city"`
	Country       string             `json:"country" validate:"required"`
	Note          *string            `json:"note"`
	PlaceID       *string            `json:"placeId"`
	TripAdvisorID *string            `json:"tripAdvisorId"`
	FourSquareID  *string            `json:"fourSquareId"`
	Tags          []string           `json:"tags"`
	StartDateTime time.Time          `json:"startDateTime"`
	EndDateTime   time.Time          `json:"endDateTime"`
}
