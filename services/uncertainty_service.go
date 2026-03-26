package services

import (
	"battery-agent/configs"
	"battery-agent/models"
	"context"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Thresholds defines normal operating ranges
var Thresholds = map[string]struct {
	Min      float64
	Max      float64
	Critical float64
}{
	"soc":         {Min: 10, Max: 100, Critical: 5},
	"soh":         {Min: 60, Max: 100, Critical: 40},
	"voltage":     {Min: 42, Max: 54.6, Critical: 40},
	"current":     {Min: -30, Max: 30, Critical: 40},
	"temperature": {Min: -10, Max: 45, Critical: 55},
}

// CheckBatteryUncertainties analyzes a battery for anomalies
func CheckBatteryUncertainties(imei string) ([]models.UncertaintyRecord, error) {
	metrics, err := FetchFullSnapshot(imei)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics for %s: %w", imei, err)
	}

	var uncertainties []models.UncertaintyRecord

	for metricName, value := range metrics {
		threshold, exists := Thresholds[metricName]
		if !exists {
			continue
		}

		var deviation float64
		var severity string

		if value < threshold.Min {
			deviation = threshold.Min - value
			severity = classifySeverity(deviation, threshold.Min, threshold.Critical)
		} else if value > threshold.Max {
			deviation = value - threshold.Max
			severity = classifySeverity(deviation, threshold.Max, threshold.Critical)
		} else {
			continue
		}

		expected := (threshold.Min + threshold.Max) / 2

		uncertainties = append(uncertainties, models.UncertaintyRecord{
			IMEI:       imei,
			MetricName: metricName,
			Expected:   expected,
			Actual:     value,
			Deviation:  deviation,
			Severity:   severity,
			Resolved:   false,
			Timestamp:  time.Now(),
		})
	}

	// Store in MongoDB
	if len(uncertainties) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var docs []interface{}
		for _, u := range uncertainties {
			docs = append(docs, u)
		}
		configs.UncertaintyCollection.InsertMany(ctx, docs)
	}

	return uncertainties, nil
}

// CheckSOHDegradation analyzes SOH trend over 30 days
func CheckSOHDegradation(imei string) (*models.UncertaintyRecord, error) {
	points, err := FetchMetricHistory(imei, "soh", 720) // 30 days
	if err != nil || len(points) < 2 {
		return nil, fmt.Errorf("insufficient SOH data")
	}

	firstSOH := points[0].Value
	lastSOH := points[len(points)-1].Value
	daysDiff := float64(len(points)) / 24.0
	degradationRate := (firstSOH - lastSOH) / daysDiff

	if degradationRate > 0.05 {
		severity := "medium"
		if degradationRate > 0.1 {
			severity = "high"
		}
		if degradationRate > 0.2 {
			severity = "critical"
		}

		record := &models.UncertaintyRecord{
			IMEI:       imei,
			MetricName: "soh_degradation_rate",
			Expected:   0.01,
			Actual:     degradationRate,
			Deviation:  degradationRate - 0.01,
			Severity:   severity,
			Resolved:   false,
			Timestamp:  time.Now(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		configs.UncertaintyCollection.InsertOne(ctx, record)

		return record, nil
	}

	return nil, nil
}

// DetectVoltageSpikes checks for sudden voltage anomalies
func DetectVoltageSpikes(imei string, hours int) ([]models.UncertaintyRecord, error) {
	points, err := FetchMetricHistory(imei, "voltage", hours)
	if err != nil || len(points) < 6 {
		return nil, err
	}

	var anomalies []models.UncertaintyRecord
	windowSize := 5

	for i := windowSize; i < len(points); i++ {
		window := points[i-windowSize : i]

		var sum float64
		for _, p := range window {
			sum += p.Value
		}
		mean := sum / float64(windowSize)

		var sqDiffSum float64
		for _, p := range window {
			sqDiffSum += math.Pow(p.Value-mean, 2)
		}
		stdDev := math.Sqrt(sqDiffSum / float64(windowSize))

		currentValue := points[i].Value
		deviation := math.Abs(currentValue - mean)

		if deviation > 2*stdDev && stdDev > 0.1 {
			severity := "low"
			if deviation > 3*stdDev {
				severity = "high"
			}

			anomalies = append(anomalies, models.UncertaintyRecord{
				IMEI:       imei,
				MetricName: "voltage_spike",
				Expected:   mean,
				Actual:     currentValue,
				Deviation:  deviation,
				Severity:   severity,
				Resolved:   false,
				Timestamp:  points[i].Timestamp,
			})
		}
	}

	return anomalies, nil
}

// GetUnresolvedUncertainties returns all unresolved issues for a battery
func GetUnresolvedUncertainties(imei string) ([]models.UncertaintyRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var records []models.UncertaintyRecord
	cursor, err := configs.UncertaintyCollection.Find(ctx, bson.M{
		"imei":     imei,
		"resolved": false,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	return records, nil
}

func classifySeverity(deviation, normalBound, criticalBound float64) string {
	ratio := deviation / math.Abs(criticalBound-normalBound)
	switch {
	case ratio > 0.8:
		return "critical"
	case ratio > 0.5:
		return "high"
	case ratio > 0.2:
		return "medium"
	default:
		return "low"
	}
}
