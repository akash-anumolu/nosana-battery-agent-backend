package routes

import (
	"battery-agent/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterAlertRoutes registers all alert endpoints
func RegisterAlertRoutes(api *echo.Group) {
	alerts := api.Group("/alerts")

	alerts.GET("", controllers.GetAllAlerts)
	alerts.PUT("/:id/acknowledge", controllers.AcknowledgeAlert)
}
