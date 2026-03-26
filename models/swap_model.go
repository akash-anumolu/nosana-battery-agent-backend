package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SwapEvent records each battery swap transaction
type SwapEvent struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StationID      string             `json:"station_id" bson:"station_id"`
	BatteryOutIMEI string             `json:"battery_out_imei" bson:"battery_out_imei"`
	BatteryInIMEI  string             `json:"battery_in_imei" bson:"battery_in_imei"`
	SOCOut         float64            `json:"soc_out" bson:"soc_out"`
	SOCIn          float64            `json:"soc_in" bson:"soc_in"`
	UserId         primitive.ObjectID `json:"user_id" bson:"user_id"`
	Duration       float64            `json:"duration_seconds" bson:"duration_seconds"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
}
