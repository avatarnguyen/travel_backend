package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Iternary struct {
	ID           primitive.ObjectID `bson:"_id"`
	IternaryID   string             `json:"iternaryId"`
	Name         *string            `json:"name"`
	Description  *string            `json:"description"`
	Note         *string            `json:"note"`
	City         *string            `json:"city"`
	Country      *string            `json:"country"`
	Creator      string             `json:"creator"`
	Owner        *string            `json:"owner"`
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
	Note          *string            `json:"note"`
	Place         *Place             `json:"place"`
	StartDateTime time.Time          `json:"startDateTime"`
	EndDateTime   time.Time          `json:"endDateTime"`
}
