package routes

import (
	"battery-agent/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterStationRoutes registers all station endpoints
func RegisterStationRoutes(api *echo.Group) {
	stations := api.Group("/stations")

	stations.GET("", controllers.GetAllStations)
	stations.GET("/:id/dashboard", controllers.GetStationDashboard)
}
