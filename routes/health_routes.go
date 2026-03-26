package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterHealthRoutes registers health check endpoint
func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "battery-agent",
		})
	})
}
