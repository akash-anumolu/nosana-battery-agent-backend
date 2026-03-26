package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LocationLog represents geographical and odometer data
type LocationLog struct {
	Lat float64 `json:"lat,omitempty" bson:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty" bson:"lng,omitempty"`
	Odo float64 `json:"odo,omitempty" bson:"odo,omitempty"`
}

// BatteryDetails represents static battery specifications
type BatteryDetails struct {
	BatteryDateOfManufacture string `json:"battery_date_of_manufacture" bson:"battery_date_of_manufacture,omitempty"`
	BatteryManufacturer      string `json:"battery_manufacturer" bson:"battery_manufacturer,omitempty"`
	BatteryModelName         string `json:"battery_model_name" bson:"battery_model_name,omitempty"`
	BatteryType              string `json:"battery_type" bson:"battery_type,omitempty"`
	BatteryVoltage           string `json:"battery_voltage" bson:"battery_voltage,omitempty"`
	BatteryCapacity          string `json:"battery_capacity" bson:"battery_capacity,omitempty"`
	BatteryAdditionalInfo    string `json:"battery_additional_info" bson:"battery_additional_info,omitempty"`
	BatteryStatus            string `json:"battery_status,omitempty" bson:"battery_status,omitempty"`
}

// BatteryLog represents live battery state from IoT device
type BatteryLog struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IMEI              string             `json:"imei,omitempty" bson:"imei,omitempty"`
	SerialNumber      string             `json:"serialnumber,omitempty" bson:"serialnumber,omitempty"`
	CreatedAt         time.Time          `json:"createdtime,omitempty" bson:"createdtime,omitempty"`
	UpdatedAt         time.Time          `json:"updatedtime,omitempty" bson:"updatedtime,omitempty"`
	BatteryDetails    BatteryDetails     `json:"batterydetails,omitempty" bson:"batterydetails,omitempty"`
	ConnectedSystem   string             `json:"connectedsystem,omitempty" bson:"connectedsystem,omitempty"`
	Geotag            LocationLog        `json:"geotag,omitempty" bson:"locationlog,omitempty"`
	SOC               float64            `json:"soc,omitempty" bson:"soc,omitempty"`
	SOH               float64            `json:"soh,omitempty" bson:"soh,omitempty"`
	Current           float64            `json:"current,omitempty" bson:"current,omitempty"`
	CycleCount        int                `json:"cyclecount" bson:"cyclecount"`
	Voltage           float64            `json:"voltage,omitempty" bson:"voltage,omitempty"`
	ChargingStatus    int                `json:"chargingStatus,omitempty" bson:"chargingStatus,omitempty"`
	SetChargingStatus *int               `json:"setchargingstatus,omitempty" bson:"setchargingstatus,omitempty"`
	UserId            primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
}
