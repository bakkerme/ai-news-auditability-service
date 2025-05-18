# AI News Auditability Service

This service provides observability and quality measurement for AI-generated content from the AI News Processor. It allows tracking LLM outputs over time, evaluating quality metrics, and providing insights for optimization.

## API Overview

The service exposes a RESTful API with these main capabilities:

1. **Run Data Management**
   - Submit processing run data from the AI News Processor
   - Retrieve run data and metadata
   - List available processing runs

2. **Benchmark Evaluation**
   - Trigger benchmark evaluations of processing runs
   - Retrieve benchmark results
   - Compare quality metrics across runs

3. **Metrics & Analysis**
   - View metrics by persona
   - Track quality metrics over time
   - Generate insights for optimization

## Implementation Path

### Phase 1: Core API and Storage
- [x] Define OpenAPI specification
- [ ] Set up Go project structure
- [ ] Implement Echo framework setup
- [ ] Create BadgerDB connection and configuration
- [ ] Implement run data ingestion endpoint
- [ ] Implement basic data storage and retrieval

### Phase 2: Benchmark Engine
- [ ] Integrate existing benchmark logic from the processor's benchmark tool
- [ ] Build benchmark queue system
- [ ] Implement benchmark triggers and processing
- [ ] Store and retrieve benchmark results

### Phase 3: Metrics and Analytics
- [ ] Implement metrics calculation from run and benchmark data
- [ ] Create historical data queries and aggregations
- [ ] Build time-series metrics endpoints
- [ ] Add persona-specific metrics

### Phase 4: UI Development
- [ ] Design and implement basic dashboard
- [ ] Add run data viewing
- [ ] Create visualizations for quality metrics
- [ ] Implement comparison views

## Getting Started

For now, review the OpenAPI specification in `planning/api-doc.yaml` to understand the API structure and data models.

As development progresses, we'll add instructions for running the service locally, configuration options, and deployment guidance. 