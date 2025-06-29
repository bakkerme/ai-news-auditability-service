package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/bakkerme/ai-news-auditability-service/internal/customerrors"
	"github.com/bakkerme/ai-news-auditability-service/internal/models"
	"github.com/bakkerme/ai-news-auditability-service/internal/openai"
	"github.com/bakkerme/ai-news-auditability-service/internal/storage"
	anpmodels "github.com/bakkerme/ai-news-processor/models"
	"github.com/google/uuid"
)

const evaluationPrompt = `You are an expert in evaluating AI-generated content. Your task is to evaluate the quality of the following post summary, focusing purely on how well it summarizes and analyzes the content.

The persona is {{.PersonaIdentity}}

The persona's focus areas are:
{{range .FocusAreas}}* {{.}}
{{end}}

The summary should be marked as irrelevant if it matches:
{{range .ExclusionCriteria}}* {{.}}
{{end}}

For each summary, evaluate how well it summarizes the post, focusing on the following criteria:

1. Summary Quality (choose one):
   - Excellent: Comprehensive summary that captures all key details and provides a clear, well-structured overview
   - Good: Clear summary with some details but lacks depth or clarity
   - Fair: Basic summary with some details but lacks depth or clarity
   - Poor: Incomplete or unclear summary lacking essential details

2. Evaluation Criteria:
   - Comprehensiveness: Does it capture all key details?
   - Technical Accuracy: If technical details are provided, are they accurate?
   - Clarity: Is the information presented in a clear, well-structured manner?
   - Comment Integration: Are community discussions and feedback well-analyzed?

3. Relevance Assessment (separate from quality rating):
   - Check if the original content matches any exclusion criteria. If it does, the IsRelevant flag should be false.
   - Evaluate if the IsRelevant flag is set appropriately
   - Assess if the relevance explanation is clear and justified

Respond with a JSON object containing:
{
  "quality_rating": string,  // One of: "Excellent", "Good", "Fair", "Poor"
  "quality_explanation": string,  // Detailed explanation of the summary quality
  "relevance_correct": boolean,  // Whether IsRelevant flag was set correctly based on exclusion criteria
  "relevance_explanation": string // Explanation of relevance assessment
}`

// BenchmarkService handles benchmark operations
type BenchmarkService struct {
	llmURL    string
	llmAPIKey string
	llmModel  string
}

// NewBenchmarkService creates a new benchmark service
func NewBenchmarkService(llmURL, llmAPIKey, llmModel string) *BenchmarkService {
	return &BenchmarkService{
		llmURL:    llmURL,
		llmAPIKey: llmAPIKey,
		llmModel:  llmModel,
	}
}

// EvaluationResult represents the structure of the benchmark evaluation response
type EvaluationResult struct {
	QualityRating        string `json:"quality_rating"`
	QualityExplanation   string `json:"quality_explanation"`
	RelevanceExplanation string `json:"relevance_explanation"`
	RelevanceCorrect     bool   `json:"relevance_correct"`
}

// EvaluationResultSchema defines the JSON schema for the evaluation result
var EvaluationResultSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"quality_rating": map[string]interface{}{
			"type":        "string",
			"description": "Descriptive rating for summary quality",
			"enum":        []string{"Excellent", "Good", "Fair", "Poor"},
		},
		"quality_explanation": map[string]interface{}{
			"type":        "string",
			"description": "Detailed explanation of the rating",
		},
		"relevance_explanation": map[string]interface{}{
			"type":        "string",
			"description": "Explanation of relevance assessment",
		},
		"relevance_correct": map[string]interface{}{
			"type":        "boolean",
			"description": "Whether IsRelevant flag was set correctly",
		},
	},
	"required": []string{"quality_rating", "quality_explanation", "relevance_explanation", "relevance_correct"},
	"additionalProperties": false,
}

// CreateBenchmark creates a new benchmark for the given run ID
func (bs *BenchmarkService) CreateBenchmark(runID string) (*models.BenchmarkResponse, error) {
	// Check if the run exists
	runData, err := storage.GetRunData(runID)
	if err != nil {
		return nil, fmt.Errorf("failed to get run data: %w", err)
	}

	// Generate a unique benchmark ID
	benchmarkID := uuid.NewString()

	// Create benchmark response
	response := &models.BenchmarkResponse{
		ID:                      benchmarkID,
		Status:                  "queued",
		Message:                 "Benchmark queued for processing",
		EstimatedCompletionTime: time.Now().Add(5 * time.Minute), // Rough estimate
	}

	// Start benchmark processing in a goroutine
	go bs.processBenchmark(benchmarkID, runID, runData)

	return response, nil
}

// processBenchmark performs the actual benchmark evaluation
func (bs *BenchmarkService) processBenchmark(benchmarkID, runID string, runData *models.PersistedRunData) {
	log.Printf("Starting benchmark processing for run ID: %s, benchmark ID: %s", runID, benchmarkID)

	// Initialize OpenAI client
	llmClient := openai.New(bs.llmURL, bs.llmAPIKey, bs.llmModel)

	// Generate evaluation prompt with persona-specific information
	tmpl, err := template.New("evaluation").Parse(evaluationPrompt)
	if err != nil {
		log.Printf("Error parsing evaluation prompt template: %v", err)
		bs.saveBenchmarkError(benchmarkID, runID, "Failed to parse evaluation prompt", err)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, runData.Persona)
	if err != nil {
		log.Printf("Error executing evaluation prompt template: %v", err)
		bs.saveBenchmarkError(benchmarkID, runID, "Failed to execute evaluation prompt", err)
		return
	}

	fullPrompt := buf.String()

	// Prepare benchmark results
	results := &models.BenchmarkResults{
		BenchmarkID:         benchmarkID,
		RunID:               runID,
		PersonaName:         runData.Persona.Name,
		PersonaFocusAreas:   runData.Persona.FocusAreas,
		DetailedEvaluations: make(map[string]models.EvaluationResult),
		MissingItems:        make([]string, 0),
		Timestamp:           time.Now(),
	}

	// Build a map from ID to raw input for matching
	rawInputByID := make(map[string]string)
	processedIDs := make(map[string]bool)

	// Extract IDs from the raw input in overall summaries
	for _, summary := range runData.EntrySummaries {
		// Try to extract the ID from the raw input (assuming 'ID: <id>' is present)
		lines := strings.Split(summary.RawInput, "\n")
		var id string
		for _, line := range lines {
			if strings.HasPrefix(line, "ID: ") {
				id = strings.TrimSpace(strings.TrimPrefix(line, "ID: "))
				break
			}
		}
		if id != "" {
			rawInputByID[id] = summary.RawInput
		}
	}

	// Process each item in the benchmark data
	for _, result := range runData.EntrySummaries {
		if result.Results.ID == "" {
			log.Printf("Warning: Empty ID for result")
			continue
		}

		processedIDs[result.Results.ID] = true
		log.Printf("Processing entry (ID: %s)...", result.Results.ID)

		// Find the matching raw input by ID
		rawInput, ok := rawInputByID[result.Results.ID]
		if !ok {
			log.Printf("Warning: No matching raw input for result ID: %s", result.Results.ID)
			continue
		}

		// Create evaluation input
		evaluationInput := fmt.Sprintf("Source Material:\n%s\n\nGenerated Summary:\n%s\n",
			rawInput,
			bs.formatSummary(result.Results))

		// Call LLM for evaluation
		log.Printf("Calling LLM for evaluation of entry ID: %s...", result.Results.ID)
		resultChan := make(chan customerrors.ErrorString, 1)
		bs.chatCompletionForBenchmarkEvaluation(llmClient, fullPrompt, []string{evaluationInput}, resultChan)
		evalResponse := <-resultChan
		if evalResponse.Err != nil {
			log.Printf("Error evaluating entry %s: %v", result.Results.ID, evalResponse.Err)
			continue
		}

		// Parse evaluation result
		var evalResult EvaluationResult
		jsonStr := llmClient.PreprocessJSON(evalResponse.Value)
		err = json.Unmarshal([]byte(jsonStr), &evalResult)
		if err != nil {
			log.Printf("Error parsing evaluation result for %s: %v", result.Results.ID, err)
			continue
		}

		log.Printf("Evaluation for entry ID %s: Quality Rating = %s, Relevance Correct = %v",
			result.Results.ID, evalResult.QualityRating, evalResult.RelevanceCorrect)

		// Convert to models.EvaluationResult
		modelEvalResult := models.EvaluationResult{
			QualityRating:        evalResult.QualityRating,
			QualityExplanation:   evalResult.QualityExplanation,
			RelevanceCorrect:     evalResult.RelevanceCorrect,
			RelevanceExplanation: evalResult.RelevanceExplanation,
		}

		results.DetailedEvaluations[result.Results.ID] = modelEvalResult
		results.TotalItems++
	}

	// Check for missing items
	for id := range rawInputByID {
		if !processedIDs[id] {
			log.Printf("Found missing item (ID: %s)...", id)
			results.MissingItems = append(results.MissingItems, id)

			// Add a Poor rating evaluation for the missing item
			results.DetailedEvaluations[id] = models.EvaluationResult{
				QualityRating:        "Poor",
				QualityExplanation:   "Item was present in raw input but missing from processed results",
				RelevanceCorrect:     false,
				RelevanceExplanation: "Unable to assess relevance as item was not processed",
			}
			results.TotalItems++
		}
	}

	// Calculate aggregate metrics
	log.Println("Calculating aggregate metrics...")
	var correctRelevance int
	for _, eval := range results.DetailedEvaluations {
		if eval.RelevanceCorrect {
			correctRelevance++
		}
	}

	if results.TotalItems > 0 {
		results.RelevanceAccuracy = float64(correctRelevance) / float64(results.TotalItems)

		// Calculate quality score with Poor rated at 0%
		var totalQualityScore float64
		for _, eval := range results.DetailedEvaluations {
			switch eval.QualityRating {
			case "Excellent":
				totalQualityScore += 100.0
			case "Good":
				totalQualityScore += 75.0
			case "Fair":
				totalQualityScore += 50.0
			case "Poor":
				totalQualityScore += 0.0
			}
		}
		results.QualityScore = totalQualityScore / float64(results.TotalItems)
	}

	// Save benchmark results
	err = bs.saveBenchmarkResults(benchmarkID, results)
	if err != nil {
		log.Printf("Error saving benchmark results: %v", err)
		bs.saveBenchmarkError(benchmarkID, runID, "Failed to save benchmark results", err)
		return
	}

	log.Printf("Benchmark processing completed for run ID: %s, benchmark ID: %s", runID, benchmarkID)
}

// chatCompletionForBenchmarkEvaluation queries the LLM for a benchmark evaluation
func (bs *BenchmarkService) chatCompletionForBenchmarkEvaluation(llmClient openai.OpenAIClient, systemPrompt string, userPrompts []string, results chan customerrors.ErrorString) {
	schemaParams := &openai.SchemaParameters{
		Schema:      EvaluationResultSchema,
		Name:        "benchmark_evaluation",
		Description: "an object representing a benchmark evaluation result (quality and relevance)",
	}

	// Setting temperature to 0.0 for more consistent evaluations
	temperature := 0.0

	llmClient.ChatCompletion(
		systemPrompt,
		userPrompts,
		[]string{}, // No image URLs
		schemaParams,
		temperature,
		0, // No max tokens limit
		results,
	)
}

// formatSummary formats an Item into a string for evaluation
func (bs *BenchmarkService) formatSummary(item anpmodels.Item) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Title: %s\n", item.Title))
	summary.WriteString(fmt.Sprintf("ID: %s\n", item.ID))
	summary.WriteString(fmt.Sprintf("Summary: %s\n", item.Summary))
	summary.WriteString(fmt.Sprintf("Comment Summary: %s\n", item.CommentSummary))
	summary.WriteString(fmt.Sprintf("Image Summary: %s\n", item.ImageSummary))
	summary.WriteString(fmt.Sprintf("Web Summary: %s\n", item.WebContentSummary))
	summary.WriteString(fmt.Sprintf("IsRelevant: %v\n", item.IsRelevant))
	return summary.String()
}

// saveBenchmarkResults saves benchmark results to storage
func (bs *BenchmarkService) saveBenchmarkResults(benchmarkID string, results *models.BenchmarkResults) error {
	return storage.SaveBenchmarkResults(benchmarkID, *results)
}

// saveBenchmarkError saves benchmark error information
func (bs *BenchmarkService) saveBenchmarkError(benchmarkID, runID, message string, err error) {
	errorResults := &models.BenchmarkResults{
		BenchmarkID:   benchmarkID,
		RunID:         runID,
		Timestamp:     time.Now(),
		FailureReason: fmt.Sprintf("%s: %v", message, err),
	}
	
	if saveErr := storage.SaveBenchmarkResults(benchmarkID, *errorResults); saveErr != nil {
		log.Printf("Failed to save benchmark error: %v", saveErr)
	}
}

// GetBenchmarkResults retrieves benchmark results by benchmark ID
func (bs *BenchmarkService) GetBenchmarkResults(runID string) (*models.BenchmarkResults, error) {
	return storage.GetBenchmarkResults(runID)
}