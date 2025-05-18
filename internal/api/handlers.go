package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bakkerme/ai-news-auditability-service/internal/models"
	"github.com/bakkerme/ai-news-auditability-service/internal/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SubmitRun handles POST /runs
func SubmitRun(c echo.Context) error {
	var runData models.PersistedRunData
	if err := c.Bind(&runData); err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Code: http.StatusBadRequest, Message: "Invalid run data format: " + err.Error()})
	}

	// Generate an ID for the run submission if not already present or use one from input
	submissionID := uuid.NewString()
	if runData.RunID == "" {
		runData.RunID = submissionID
	} else {
		submissionID = runData.RunID // Ensure submissionID matches the one potentially in data
	}

	if runData.RunDate.IsZero() { // Ensure RunDate is set
		runData.RunDate = time.Now()
	}

	fmt.Printf("Run data: %+v\n", runData)

	if err := storage.SaveRunData(submissionID, runData); err != nil {
		log.Printf("Error saving run data: %v", err) // Log the error
		return c.JSON(http.StatusInternalServerError, models.Error{Code: http.StatusInternalServerError, Message: "Failed to save run data: " + err.Error()})
	}

	fmt.Println("Run data received and saved for ID:", submissionID)

	response := models.RunResponse{
		ID:      submissionID,
		Message: "Run data successfully received and stored",
		Status:  "stored",
	}
	return c.JSON(http.StatusCreated, response)
}

// ListRuns handles GET /runs
func ListRuns(c echo.Context) error {
	// TODO: Implement filtering (persona, from, to) and pagination
	// For now, fetch a small number of recent runs as metadata
	runsMetadata, err := storage.ListRunMetadata(10) // Fetch latest 10
	if err != nil {
		log.Printf("Error listing runs: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Error{Code: http.StatusInternalServerError, Message: "Failed to retrieve runs: " + err.Error()})
	}

	if runsMetadata == nil { // Ensure we return an empty slice, not null, if no runs
		runsMetadata = []models.RunMetadata{}
	}

	return c.JSON(http.StatusOK, runsMetadata)
}

// GetRun handles GET /runs/{runId}
func GetRun(c echo.Context) error {
	runID := c.Param("runId")
	runData, err := storage.GetRunData(runID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") { // Corrected error check
			return c.JSON(http.StatusNotFound, models.Error{Code: http.StatusNotFound, Message: fmt.Sprintf("Run with ID '%s' not found", runID)})
		}
		log.Printf("Error getting run data for ID %s: %v", runID, err)
		return c.JSON(http.StatusInternalServerError, models.Error{Code: http.StatusInternalServerError, Message: "Failed to retrieve run data: " + err.Error()})
	}

	return c.JSON(http.StatusOK, runData)
}

// CreateBenchmark handles POST /benchmarks/create/{runId}
func CreateBenchmark(c echo.Context) error {
	runID := c.Param("runId") // runID is available for use in TODO logic
	// TODO: Implement benchmark creation logic, associate with runID

	// Dummy response for now
	if runID == "run-to-fail-benchmark-creation" { // Simulate a scenario
		return c.JSON(http.StatusNotFound, models.Error{Code: http.StatusNotFound, Message: "Run not found, cannot create benchmark"})
	}

	response := models.BenchmarkResponse{
		ID:      uuid.NewString(), // Changed from BenchmarkID, RunID removed
		Message: "Benchmark triggered successfully",
		Status:  "pending", // Timestamp removed, EstimatedCompletionTime can be added when known
	}
	return c.JSON(http.StatusAccepted, response)
}

// GetBenchmark handles GET /benchmarks/{runId}
func GetBenchmark(c echo.Context) error {
	runID := c.Param("runId")
	// TODO: Implement benchmark results retrieval logic for the given runID

	// Dummy response for now
	if runID == "nonexistent-benchmark" { // Simulate not found
		return c.JSON(http.StatusNotFound, models.Error{Code: http.StatusNotFound, Message: "Benchmark not found"})
	}

	dummyResults := models.BenchmarkResults{
		// Fields from the current model definition
		TotalItems:        10,
		RelevanceAccuracy: 0.85,
		QualityScore:      0.92,
		DetailedEvaluations: map[string]models.EvaluationResult{
			"item123": {
				QualityRating:        "Good",
				QualityExplanation:   "Summary was concise and accurate.",
				RelevanceCorrect:     true,
				RelevanceExplanation: "Item was correctly marked as relevant.",
			},
		},
		PersonaName:       "TestPersonaForBenchmark",
		PersonaFocusAreas: []string{"testing", "benchmarking"},
		MissingItems:      []string{"item456"},

		// Fields from previous simpler version (kept as they are in the Go struct)
		BenchmarkID: uuid.NewString(),               // This could be the ID of the benchmark job
		RunID:       runID,                          // The ID of the run that was benchmarked
		Timestamp:   time.Now().Add(-1 * time.Hour), // Time of results generation
		Judgement:   "The summary is good overall.",
		// Status and Scores fields removed as they are not in the current struct definition
	}
	return c.JSON(http.StatusOK, dummyResults)
}

// GetBenchmarkLogs handles GET /benchmarks/{runId}/logs
func GetBenchmarkLogs(c echo.Context) error {
	runID := c.Param("runId")
	// TODO: Implement log retrieval logic
	// from := c.QueryParam("from")
	// to := c.QueryParam("to")
	// level := c.QueryParam("level")

	// Dummy response for now
	if runID == "benchmark-with-no-logs" { // Simulate not found
		return c.JSON(http.StatusNotFound, models.Error{Code: http.StatusNotFound, Message: "Benchmark not found or no logs available"})
	}

	dummyLogs := []models.LogEntry{
		{Timestamp: time.Now().Add(-5 * time.Minute), Level: "info", Message: "Benchmark started for run: " + runID, Source: "benchmark_evaluator"},
		{Timestamp: time.Now().Add(-4 * time.Minute), Level: "debug", Message: "Loading data...", Source: "benchmark_evaluator"},
		{Timestamp: time.Now().Add(-3 * time.Minute), Level: "info", Message: "Evaluation complete.", Source: "benchmark_evaluator"},
	}
	return c.JSON(http.StatusOK, dummyLogs)
}

// StreamBenchmarkLogs handles GET /benchmarks/{runId}/logs/stream (WebSocket)
func StreamBenchmarkLogs(c echo.Context) error {
	// runID := c.Param("runId")
	// TODO: Implement WebSocket logic for streaming logs
	// This is a simplified placeholder. Real WebSocket handling is more complex.
	return c.String(http.StatusNotImplemented, "WebSocket streaming not implemented yet.")
}

// GetPersonaMetrics handles GET /metrics/persona/{personaName}
func GetPersonaMetrics(c echo.Context) error {
	// personaName := c.Param("personaName")
	// TODO: Implement metrics retrieval logic for the persona
	// from := c.QueryParam("from")
	// to := c.QueryParam("to")

	// Dummy response for now, using the placeholder PersonaMetrics model
	dummyPersonaMetrics := models.PersonaMetrics{
		Data: map[string]interface{}{
			"averageScore": 0.88,
			"runsAnalyzed": 150,
			"errorRate":    0.02,
		},
	}
	return c.JSON(http.StatusOK, dummyPersonaMetrics)
}

// GetQualityMetrics handles GET /metrics/quality
func GetQualityMetrics(c echo.Context) error {
	// TODO: Implement quality metrics retrieval logic
	// persona := c.QueryParam("persona")
	// from := c.QueryParam("from")
	// to := c.QueryParam("to")

	// Dummy response for now
	dummyQualityMetrics := models.QualityMetrics{
		Data: map[string]interface{}{
			"trend":                 "improving",
			"averageScoreLast7Days": 0.9,
		},
	}
	return c.JSON(http.StatusOK, dummyQualityMetrics)
}
