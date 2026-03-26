package routes

import (
	"battery-agent/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterBatteryRoutes registers all battery endpoints
func RegisterBatteryRoutes(api *echo.Group) {
	batteries := api.Group("/batteries")

	batteries.GET("", controllers.GetAllBatteries)
	batteries.GET("/:imei", controllers.GetBatteryByIMEI)
	batteries.GET("/:imei/live", controllers.GetBatteryLiveMetrics)
	batteries.GET("/:imei/timeseries", controllers.GetBatteryTimeSeries)
	batteries.GET("/:imei/uncertainties", controllers.GetBatteryUncertainties)
}
