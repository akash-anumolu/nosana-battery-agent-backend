package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SwapStation represents a physical battery swap station
type SwapStation struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StationID    string             `json:"station_id" bson:"station_id"`
	Name         string             `json:"name" bson:"name"`
	Location     LocationLog        `json:"location" bson:"location"`
	Status       string             `json:"status" bson:"status"`
	Capacity     int                `json:"capacity" bson:"capacity"`
	Available    int                `json:"available" bson:"available"`
	BatteryIMEIs []string           `json:"battery_imeis" bson:"battery_imeis"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
