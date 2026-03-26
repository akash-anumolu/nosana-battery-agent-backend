package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UncertaintyRecord tracks anomalies detected in battery parameters
type UncertaintyRecord struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IMEI       string             `json:"imei" bson:"imei"`
	StationID  string             `json:"station_id" bson:"station_id"`
	MetricName string             `json:"metric_name" bson:"metric_name"`
	Expected   float64            `json:"expected" bson:"expected"`
	Actual     float64            `json:"actual" bson:"actual"`
	Deviation  float64            `json:"deviation" bson:"deviation"`
	Severity   string             `json:"severity" bson:"severity"`
	Resolved   bool               `json:"resolved" bson:"resolved"`
	Timestamp  time.Time          `json:"timestamp" bson:"timestamp"`
}
