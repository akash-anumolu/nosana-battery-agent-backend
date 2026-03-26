package controllers

import (
	"battery-agent/responses"
	"battery-agent/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAllStations - GET /api/v1/stations
func GetAllStations(c echo.Context) error {
	stations, err := services.GetAllStations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching stations",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    stations,
	})
}

// GetStationDashboard - GET /api/v1/stations/:id/dashboard
func GetStationDashboard(c echo.Context) error {
	stationID := c.Param("id")

	station, err := services.GetStationByID(stationID)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.APIResponse{
			Status:  http.StatusNotFound,
			Message: "Station not found",
			Data:    stationID,
		})
	}

	batteries, _ := services.GetStationBatteries(stationID)
	swaps, _ := services.GetRecentSwaps(stationID, 20)
	alerts, _ := services.GetStationAlerts(stationID)

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Station dashboard",
		Data: map[string]interface{}{
			"station":       station,
			"batteries":     batteries,
			"battery_count": len(batteries),
			"recent_swaps":  swaps,
			"active_alerts": alerts,
			"alert_count":   len(alerts),
		},
	})
}
