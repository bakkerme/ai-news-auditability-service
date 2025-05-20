package models

import (
	"time"

	anpmodels "github.com/bakkerme/ai-news-processor/models"
)

// PersistedRunData embeds the shared RunData and adds a RunID.
// The anpmodels.RunData now uses anpmodels.AuditablePersona.
type PersistedRunData struct {
	anpmodels.RunData `json:",inline"` // Embeds fields from anpmodels.RunData
	RunID             string           `json:"runId"`
}

// RunResponse is the response after submitting run data.
// Based on #/components/schemas/RunResponse
type RunResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"` // e.g., received, processing, completed, failed
	Message string `json:"message"`
}

// Error represents a generic error response.
// Based on #/components/schemas/Error
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// RunMetadata contains basic metadata for a run for listing.
// Updated based on #/components/schemas/RunMetadata
type RunMetadata struct {
	ID               string    `json:"id"`
	PersonaName      string    `json:"personaName"`
	RunDate          time.Time `json:"runDate"` // Note: Spec says string, format: date-time. Using time.Time for Go.
	OverallModelUsed string    `json:"overallModelUsed,omitempty"`
	TotalItems       int       `json:"totalItems,omitempty"` // This would be len(RunData.EntrySummaries)
	HasBenchmark     bool      `json:"hasBenchmark,omitempty"`
}

// BenchmarkResponse is the response after triggering a benchmark.
// Updated based on #/components/schemas/BenchmarkResponse
type BenchmarkResponse struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"` // e.g. queued, processing, completed, failed
	Message                 string    `json:"message"`
	EstimatedCompletionTime time.Time `json:"estimatedCompletionTime,omitempty"` // Note: Spec says string, format: date-time.
}

// EvaluationResult holds detailed evaluation for an item.
// Based on #/components/schemas/EvaluationResult
type EvaluationResult struct {
	QualityRating        string `json:"qualityRating"` // Excellent, Good, Fair, Poor
	QualityExplanation   string `json:"qualityExplanation"`
	RelevanceCorrect     bool   `json:"relevanceCorrect"`
	RelevanceExplanation string `json:"relevanceExplanation"`
}

// BenchmarkResults contains the results of a benchmark evaluation.
// Updated based on #/components/schemas/BenchmarkResults
type BenchmarkResults struct {
	TotalItems          int                         `json:"totalItems"`
	RelevanceAccuracy   float64                     `json:"relevanceAccuracy"`
	QualityScore        float64                     `json:"qualityScore"`
	DetailedEvaluations map[string]EvaluationResult `json:"detailedEvaluations"` // Map of item ID to detailed evaluation
	PersonaName         string                      `json:"personaName"`
	PersonaFocusAreas   []string                    `json:"personaFocusAreas"`
	MissingItems        []string                    `json:"missingItems,omitempty"`
	BenchmarkID         string                      `json:"benchmarkId,omitempty"`
	RunID               string                      `json:"runId,omitempty"`
	Timestamp           time.Time                   `json:"timestamp,omitempty"`
	RawOutput           map[string]interface{}      `json:"rawOutput,omitempty"`
	Judgement           string                      `json:"judgement,omitempty"`
	FailureReason       string                      `json:"failureReason,omitempty"`
}

// LogEntry represents a single log entry.
// Updated based on #/components/schemas/LogEntry
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"` // debug, info, warn, error
	Message   string            `json:"message"`
	ItemID    string            `json:"itemId,omitempty"`
	Phase     string            `json:"phase,omitempty"` // e.g., "initialization", "evaluation"
	Progress  *LogEntryProgress `json:"progress,omitempty"`
	Type      string            `json:"type,omitempty"` // log, progress, status, complete, error (for WebSocket)
	Source    string            `json:"source,omitempty"`
}

// LogEntryProgress provides progress information within a log entry.
// Based on #/components/schemas/LogEntry/properties/progress
type LogEntryProgress struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

// BenchmarkProgress for WebSocket messages.
// Based on #/components/schemas/BenchmarkProgress
type BenchmarkProgress struct {
	Current    int     `json:"current"`
	Total      int     `json:"total"`
	Phase      string  `json:"phase"`
	Percentage float64 `json:"percentage,omitempty"`
	ItemID     string  `json:"itemId,omitempty"`
}

// BenchmarkStatus for WebSocket messages.
// Based on #/components/schemas/BenchmarkStatus
type BenchmarkStatus struct {
	Status             string    `json:"status"` // queued, initializing, processing, etc.
	Timestamp          time.Time `json:"timestamp"`
	Message            string    `json:"message,omitempty"`
	EstimatedRemaining int       `json:"estimatedRemaining,omitempty"` // in seconds
}

// BenchmarkComplete for WebSocket messages.
// Based on #/components/schemas/BenchmarkComplete
type BenchmarkComplete struct {
	RunID      string `json:"runId"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Duration   int    `json:"duration,omitempty"` // in seconds
	ResultsURL string `json:"resultsUrl,omitempty"`
}

// BenchmarkError for WebSocket messages.
// Based on #/components/schemas/BenchmarkError
type BenchmarkError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// WebSocketMessage is a generic wrapper for messages sent over WebSocket.
// Based on #/components/schemas/WebSocketMessage
type WebSocketMessage struct {
	Type string      `json:"type"` // log, progress, status, complete, error
	Data interface{} `json:"data"` // oneOf: LogEntry, BenchmarkProgress, BenchmarkStatus, BenchmarkComplete, BenchmarkError
}

// PersonaMetrics represents historical metrics for a specific persona.
// Partially defined, awaiting full spec from YAML.
// Based on #/components/schemas/PersonaMetrics
type PersonaMetrics struct {
	PersonaName string                 `json:"personaName"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// QualityMetrics represents quality metrics over time.
// Partially defined, awaiting full spec from YAML.
// Based on #/components/schemas/QualityMetrics
type QualityMetrics struct {
	Data map[string]interface{} `json:"data,omitempty"`
}
