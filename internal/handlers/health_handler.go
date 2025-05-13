package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// HealthResponse represents the response structure for health check endpoint
type HealthResponse struct {
	Status    string    `json:"status"`    // Current health status of the service
	Timestamp time.Time `json:"timestamp"` // Time when the health check was performed
	Version   string    `json:"version"`   // Current version of the application
	Uptime    string    `json:"uptime"`    // Duration since the service started
}

// HealthHandler handles health check requests and maintains service metadata
type HealthHandler struct {
	startTime time.Time // Time when the service was started
	version   string    // Version of the application
}

// NewHealthHandler creates a new instance of HealthHandler
// Parameters:
//   - version: The version of the application
//
// Returns:
//   - A pointer to the new HealthHandler instance
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthCheck handles GET requests to the health check endpoint
// It returns the current health status, version, and uptime of the service
// Parameters:
//   - c: The Fiber context containing the HTTP request and response
//
// Returns:
//   - error: Any error that occurred while processing the request
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	uptime := time.Since(h.startTime)
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
		Uptime:    uptime.String(),
	}

	log.Info().
		Str("status", response.Status).
		Str("version", response.Version).
		Str("uptime", response.Uptime).
		Msg("Health check")

	return c.JSON(response)
}
