package controllers

import (
	"battery-agent/configs"
	"battery-agent/responses"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChatRequest from frontend
type ChatRequest struct {
	Message   string   `json:"message"`
	StationID string   `json:"station_id,omitempty"`
	IMEIs     []string `json:"imeis,omitempty"`
}

// ChatResponse to frontend
type ChatResponse struct {
	Reply   string `json:"reply"`
	Context string `json:"context,omitempty"`
}

// AgentChat - POST /api/v1/agent/chat
func AgentChat(c echo.Context) error {
	var req ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, responses.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request",
			Data:    err.Error(),
		})
	}

	// Forward to ElizaOS agent
	elizaURL := fmt.Sprintf("%s/api/chat", configs.EnvElizaOSURL())

	payload, _ := json.Marshal(map[string]string{
		"message": req.Message,
	})

	resp, err := http.Post(elizaURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error contacting AI agent",
			Data:    err.Error(),
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error reading agent response",
			Data:    err.Error(),
		})
	}

	var agentReply ChatResponse
	json.Unmarshal(body, &agentReply)

	return c.JSON(http.StatusOK, responses.APIResponse{
		Status:  http.StatusOK,
		Message: "Agent response",
		Data:    agentReply,
	})
}
