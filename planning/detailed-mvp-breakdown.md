# MVP Task Breakdown

## Goal 1: Service can collect, store and display most recent run

### 1. Separate and log prompt components (DONE) 
**Subtasks:**  (ai-news-processor)
- Create struct for storing different prompt components
- Implement storage mechanism for each component (system prompt, persona, etc.)
- Build metadata collection for run context (timestamp, model used)

**Technology:**
- **Language**: Go

### 2. Build API (DONE)
**Subtasks:**
- Design RESTful API endpoints for data ingestion
- Create API documentation
- Implement API endpoints in Go
- Implement connection into ai-news-processor

**Technology:**
- **API Framework**: Echo (Go)

### 3. Build base version with data storage (Done)
**Subtasks:**
- Set up BadgerDB connection and configuration
- Implement data access layer for storage
- Create data models that map to database schema
- Set up initial data retention policies

**Technology:**
- **Platform**: Docker + Alpine Linux
- **Database**: BadgerDB

### 4. Display basic UI for data
**Subtasks:**
- Design simple dashboard for viewing most recent run
- Implement backend API endpoints for UI data
- Develop frontend components to display prompt components
- Create basic navigation structure

**Technology:**
- **Frontend Framework**: Next.js

**Detailed Plan:** See [MVP UI Planning](./ui-mvp-planning.md)

## Goal 2: Benchmark Web & Image Summarization dump acts as benchmark file

### 1. Update benchmark system for web content
**Subtasks:**
- Modify existing system to accept web content input
- Create parsers for web content formats
- Write unit tests for web content benchmarking
- Implement web content storage mechanism

### 2. Modify metrics
**Subtasks:**
- Define metrics specific to web content summarization
- Define metrics specific to image summarization
- Update calculation pipeline for new metrics
- Create comparison methods for different content types

**Technology:**
- **Evaluation**: LLM-based judge system

### 3. Decide on new image benchmark model
**Subtasks:**
- Research appropriate image benchmark models
- Document selection criteria and decision
- Create implementation plan for selected model
- Define evaluation criteria for image summarization

### 4. Build or update benchmarking system
**Subtasks:**
- Implement image benchmark model integration
- Update benchmarking pipeline for multi-content types
- Create unified scoring system across content types
- Add logging for benchmark execution

### 5. Create benchmark data parser
**Subtasks:**
- Develop script to parse benchmark data into standard format
- Implement feed data download functionality
- Add data validation for downloaded content
- Create storage system for benchmark datasets

### 6. Remove inline dump support (ai-news-processor)
**Subtasks:**
- Identify all instances of inline dump functionality
- Create migration path for existing data
- Update documentation to reflect new approach
- Validate system functionality after removal

## Goal 3: Auto Benchmarking Support

### 1. Build data storage system
**Subtasks:**
- Design schema for benchmark data and results
- Implement database structure for multiple runs
- Create data access layer for benchmark storage
- Develop compression strategy for efficient storage

**Technology:**
- **Database**: BadgerDB

### 2. Allow for multiple runs storage
**Subtasks:**
- Create unique identifiers for benchmark runs
- Implement versioning for benchmark configurations
- Build comparison functionality between runs
- Set data retention policies for benchmark history

### 3. Run benchmarks on cron or trigger
**Subtasks:**
- Implement cron job system for scheduled benchmarking
- Create trigger mechanism for on-demand benchmarking
- Add logging and monitoring for benchmark jobs
- Develop notification system for benchmark completion

**Technology:**
- **Scheduling**: Cron implementation

### 4. Provide output in Web UI
**Subtasks:**
- Design dashboard for viewing benchmark results
- Implement data visualization components for metrics
- Create comparison views for multiple runs
- Add filtering and sorting capabilities

**Technology:**
- **Visualization Library**: Recharts or Chart.js

## Goal 4: Optimization Mode

### 1. Allow service to modify Persona + Prompt
**Subtasks:**
- Create interface for modifying components
- Implement version control for modifications
- Develop audit trail for tracking changes
- Build safety mechanisms to prevent destructive changes

### 2. Use benchmark data for optimization
**Subtasks:**
- Develop optimization algorithm based on benchmark results
- Create validation pipeline across multiple benchmarks
- Implement comparison between original and optimized versions
- Build automated optimization workflow

### 3. Add rudimentary quality tracking chart
**Subtasks:**
- Design time-series charts for quality metrics
- Implement backend API for historical quality data
- Create frontend visualization components
- Add basic trend analysis functionality

**Technology:**
- **Visualization**: Simple chart implementation

## Goal 5: Deep Metrics Visualization

### 1. Track prompt version, model, persona over time
**Subtasks:**
- Create data structures for tracking component versions
- Develop model tracking functionality
- Implement relationship mapping between components
- Build historical query capabilities

### 2. Categorize entries based on quality and sources
**Subtasks:**
- Define quality categories and thresholds
- Implement automatic categorization based on metrics
- Create visualization for quality distribution
- Build data source tracking and analysis system
