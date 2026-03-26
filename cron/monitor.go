package cron

import (
	"battery-agent/models"
	"battery-agent/services"
	"fmt"
	"log"
	"time"
)

// StartMonitoringCron starts all background jobs
func StartMonitoringCron() {
	// Every 5 minutes: check all batteries
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			log.Println("🔋 Running battery monitoring cycle...")
			monitorAllBatteries()
		}
	}()

	// Every 6 hours: check SOH degradation
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		for range ticker.C {
			log.Println("📉 Running SOH degradation check...")
			checkAllSOHDegradation()
		}
	}()

	log.Println("✅ Monitoring cron jobs started")
}

func monitorAllBatteries() {
	batteries, err := services.GetAllBatteries()
	if err != nil {
		log.Printf("❌ Error fetching batteries: %v", err)
		return
	}

	for _, battery := range batteries {
		if battery.IMEI == "" {
			continue
		}

		// Check uncertainties
		uncertainties, err := services.CheckBatteryUncertainties(battery.IMEI)
		if err != nil {
			log.Printf("⚠️ Error checking %s: %v", battery.IMEI, err)
			continue
		}

		// Generate alerts for critical issues
		for _, u := range uncertainties {
			if u.Severity == "critical" || u.Severity == "high" {
				services.CreateAlert(models.Alert{
					StationID: battery.ConnectedSystem,
					IMEI:      battery.IMEI,
					Type:      "anomaly",
					Severity:  u.Severity,
					Message: fmt.Sprintf(
						"%s anomaly: expected %.2f, got %.2f (deviation: %.2f)",
						u.MetricName, u.Expected, u.Actual, u.Deviation,
					),
				})
				log.Printf("🚨 Alert: %s for %s", u.Severity, battery.IMEI)
			}
		}

		// Store time series snapshot
		services.StoreTimeSeriesSnapshot(battery.IMEI, battery.ConnectedSystem)
	}
}

func checkAllSOHDegradation() {
	batteries, err := services.GetAllBatteries()
	if err != nil {
		return
	}

	for _, battery := range batteries {
		if battery.IMEI == "" {
			continue
		}

		record, err := services.CheckSOHDegradation(battery.IMEI)
		if err != nil || record == nil {
			continue
		}

		if record.Severity == "high" || record.Severity == "critical" {
			services.CreateAlert(models.Alert{
				IMEI:     battery.IMEI,
				Type:     "degradation",
				Severity: record.Severity,
				Message: fmt.Sprintf(
					"SOH degrading at %.4f%%/day (normal: 0.01%%/day)",
					record.Actual,
				),
			})
		}
	}
}
