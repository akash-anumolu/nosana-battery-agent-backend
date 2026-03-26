package controllers

import (
	"battery-agent/responses"
	"battery-agent/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllBatteries - GET /api/v1/batteries
func GetAllBatteries(c echo.Context) error {
	batteries, err := services.GetAllBatteries()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching batteries",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    batteries,
	})
}

// GetBatteryByIMEI - GET /api/v1/batteries/:imei
func GetBatteryByIMEI(c echo.Context) error {
	imei := c.Param("imei")

	battery, err := services.GetBatteryByIMEI(imei)
	if err != nil {
		return c.JSON(http.StatusNotFound, responses.APIResponse{
			Status:  http.StatusNotFound,
			Message: "Battery not found",
			Data:    imei,
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    battery,
	})
}

// GetBatteryLiveMetrics - GET /api/v1/batteries/:imei/live
func GetBatteryLiveMetrics(c echo.Context) error {
	imei := c.Param("imei")

	metrics, err := services.FetchFullSnapshot(imei)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching live metrics",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Live metrics from Prometheus",
		Data:    metrics,
	})
}

// GetBatteryTimeSeries - GET /api/v1/batteries/:imei/timeseries?metric=soc&hours=24
func GetBatteryTimeSeries(c echo.Context) error {
	imei := c.Param("imei")
	metric := c.QueryParam("metric")
	hoursStr := c.QueryParam("hours")

	if metric == "" {
		metric = "soc"
	}

	hours := 24
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil {
			hours = h
		}
	}

	points, err := services.FetchMetricHistory(imei, metric, hours)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching time series",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Time series data",
		Data: map[string]interface{}{
			"imei":   imei,
			"metric": metric,
			"hours":  hours,
			"points": points,
			"count":  len(points),
		},
	})
}

// GetBatteryUncertainties - GET /api/v1/batteries/:imei/uncertainties
func GetBatteryUncertainties(c echo.Context) error {
	imei := c.Param("imei")

	uncertainties, err := services.CheckBatteryUncertainties(imei)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error checking uncertainties",
			Data:    err.Error(),
		})
	}

	degradation, _ := services.CheckSOHDegradation(imei)

	spikes, _ := services.DetectVoltageSpikes(imei, 24)

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Uncertainty analysis complete",
		Data: map[string]interface{}{
			"imei":            imei,
			"uncertainties":   uncertainties,
			"soh_degradation": degradation,
			"voltage_spikes":  spikes,
			"total_issues":    len(uncertainties),
		},
	})
}
