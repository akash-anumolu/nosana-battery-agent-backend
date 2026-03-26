package services

import (
	"battery-agent/configs"
	"battery-agent/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAlert stores a new alert
func CreateAlert(alert models.Alert) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alert.CreatedAt = time.Now()
	alert.Acknowledged = false

	_, err := configs.AlertCollection.InsertOne(ctx, alert)
	return err
}

// GetAlerts returns filtered alerts
func GetAlerts(severity string, acknowledged string) ([]models.Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if severity != "" {
		filter["severity"] = severity
	}
	if acknowledged != "" {
		filter["acknowledged"] = acknowledged == "true"
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(100)

	var alerts []models.Alert
	cursor, err := configs.AlertCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

// GetStationAlerts returns unresolved alerts for a station
func GetStationAlerts(stationID string) ([]models.Alert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alerts []models.Alert
	cursor, err := configs.AlertCollection.Find(ctx, bson.M{
		"station_id":   stationID,
		"acknowledged": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

// AcknowledgeAlert marks an alert as acknowledged
func AcknowledgeAlert(alertID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(alertID)
	if err != nil {
		return fmt.Errorf("invalid alert ID: %w", err)
	}

	_, err = configs.AlertCollection.UpdateByID(ctx, objID, bson.M{
		"$set": bson.M{"acknowledged": true},
	})

	return err
}
