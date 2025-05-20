package api

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all API routes with the Echo server.
func RegisterRoutes(e *echo.Echo, apiHandler *API) {
	// Base path v1
	v1 := e.Group("/v1")

	// Run Endpoints
	v1.POST("/runs", apiHandler.SubmitRun)          // Submit run data
	v1.GET("/runs", apiHandler.ListRuns)            // List runs
	v1.GET("/runs/latest", apiHandler.GetLatestRun) // Get latest run data
	v1.GET("/runs/:runId", apiHandler.GetRun)       // Get specific run data

	// Benchmark Endpoints
	v1.POST("/benchmarks/create/:runId", apiHandler.CreateBenchmark) // Create new benchmark
	v1.GET("/benchmarks/:runId", apiHandler.GetBenchmark)            // Get benchmark results
	v1.GET("/benchmarks/:runId/logs", apiHandler.GetBenchmarkLogs)   // Get benchmark logs
	// WebSocket endpoint for streaming logs
	v1.GET("/benchmarks/:runId/logs/stream", apiHandler.StreamBenchmarkLogs)

	// Metrics Endpoints
	// Note: The OpenAPI spec shows /metrics/persona/{personaName} and then other /metrics/ endpoints.
	// I'll need to check the rest of the spec for other metric endpoints.
	v1.GET("/metrics/persona/:personaName", apiHandler.GetPersonaMetrics) // Get metrics by persona
	v1.GET("/metrics/quality", apiHandler.GetQualityMetrics)              // Get quality metrics over time

	// TODO: Add other metric endpoints as defined in the full api-doc.yaml
	// e.g., /metrics/model/{modelName}, /metrics/summary, etc.
}
