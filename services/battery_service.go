package services

import (
	"battery-agent/configs"
	"battery-agent/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// FetchBatteryMetric fetches a single metric from Prometheus
func FetchBatteryMetric(imei string, metric string) (float64, error) {
	query := fmt.Sprintf(`battery_%s{imei="%s"}`, metric, imei)
	resp, err := QueryInstant(query)
	if err != nil {
		return 0, err
	}
	if len(resp.Data.Result) == 0 {
		return 0, fmt.Errorf("no %s data for IMEI: %s", metric, imei)
	}
	return ParsePromValue(resp.Data.Result[0].Value[1])
}

// FetchFullSnapshot gets all metrics for a battery
func FetchFullSnapshot(imei string) (map[string]float64, error) {
	metrics := map[string]float64{}
	metricNames := []string{"soc", "soh", "voltage", "current", "temperature"}

	for _, name := range metricNames {
		value, err := FetchBatteryMetric(imei, name)
		if err == nil {
			metrics[name] = value
		}
	}

	return metrics, nil
}

// FetchMetricHistory gets time series for any metric
func FetchMetricHistory(imei string, metric string, hours int) ([]TimeSeriesPoint, error) {
	query := fmt.Sprintf(`battery_%s{imei="%s"}`, metric, imei)
	end := time.Now()
	start := end.Add(-time.Duration(hours) * time.Hour)

	step := calculateStep(hours)

	resp, err := QueryRange(query, start, end, step)
	if err != nil {
		return nil, err
	}

	return ParseTimeSeries(resp)
}

// StoreTimeSeriesSnapshot saves current metrics to MongoDB
func StoreTimeSeriesSnapshot(imei string, stationID string) error {
	metrics, err := FetchFullSnapshot(imei)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var docs []interface{}
	now := time.Now()

	for metric, value := range metrics {
		docs = append(docs, models.BatteryTimeSeries{
			IMEI:      imei,
			StationID: stationID,
			Metric:    metric,
			Value:     value,
			Timestamp: now,
		})
	}

	if len(docs) > 0 {
		_, err = configs.TimeSeriesCollection.InsertMany(ctx, docs)
	}

	return err
}

// GetAllBatteries returns all batteries from MongoDB
func GetAllBatteries() ([]models.BatteryLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var batteries []models.BatteryLog
	cursor, err := configs.BatteryCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &batteries); err != nil {
		return nil, err
	}

	return batteries, nil
}

// GetBatteryByIMEI returns a single battery
func GetBatteryByIMEI(imei string) (*models.BatteryLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var battery models.BatteryLog
	err := configs.BatteryCollection.FindOne(ctx, bson.M{"imei": imei}).Decode(&battery)
	if err != nil {
		return nil, err
	}

	return &battery, nil
}

// helper: calculate appropriate step based on time range
func calculateStep(hours int) string {
	switch {
	case hours <= 6:
		return "1m"
	case hours <= 24:
		return "5m"
	case hours <= 168: // 1 week
		return "15m"
	default:
		return "1h"
	}
}
