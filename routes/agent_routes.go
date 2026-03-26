package routes

import (
	"battery-agent/controllers"

	"github.com/labstack/echo/v4"
)

// RegisterAgentRoutes registers AI agent chat endpoints
func RegisterAgentRoutes(api *echo.Group) {
	agent := api.Group("/agent")

	agent.POST("/chat", controllers.AgentChat)
}
