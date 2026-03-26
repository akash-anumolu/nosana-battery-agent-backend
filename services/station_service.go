package services

import (
	"battery-agent/configs"
	"battery-agent/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllStations returns all swap stations
func GetAllStations() ([]models.SwapStation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var stations []models.SwapStation
	cursor, err := configs.StationCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &stations); err != nil {
		return nil, err
	}

	return stations, nil
}

// GetStationByID returns a single station
func GetStationByID(stationID string) (*models.SwapStation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var station models.SwapStation
	err := configs.StationCollection.FindOne(ctx, bson.M{"station_id": stationID}).Decode(&station)
	if err != nil {
		return nil, err
	}

	return &station, nil
}

// GetStationBatteries returns all batteries at a station
func GetStationBatteries(stationID string) ([]models.BatteryLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var batteries []models.BatteryLog
	cursor, err := configs.BatteryCollection.Find(ctx, bson.M{"connectedsystem": stationID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &batteries); err != nil {
		return nil, err
	}

	return batteries, nil
}

// GetRecentSwaps returns recent swap events for a station
func GetRecentSwaps(stationID string, limit int64) ([]models.SwapEvent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(limit)

	var swaps []models.SwapEvent
	cursor, err := configs.SwapEventCollection.Find(ctx, bson.M{"station_id": stationID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &swaps); err != nil {
		return nil, err
	}

	return swaps, nil
}
