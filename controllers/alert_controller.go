package controllers

import (
	"battery-agent/responses"
	"battery-agent/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetAllAlerts - GET /api/v1/alerts?severity=critical&acknowledged=false
func GetAllAlerts(c echo.Context) error {
	severity := c.QueryParam("severity")
	acknowledged := c.QueryParam("acknowledged")

	alerts, err := services.GetAlerts(severity, acknowledged)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching alerts",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    alerts,
	})
}

// AcknowledgeAlert - PUT /api/v1/alerts/:id/acknowledge
func AcknowledgeAlert(c echo.Context) error {
	alertID := c.Param("id")

	err := services.AcknowledgeAlert(alertID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error acknowledging alert",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Alert acknowledged",
		Data:    alertID,
	})
}
