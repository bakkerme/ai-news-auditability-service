# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Development
```bash
# Start Go service (development)
make RunServiceDev

# Start Go service with hot reload
make RunServiceDevAir

# Start Next.js dashboard
make RunDashboardDev

# Clear BadgerDB files
make ClearBadgerDB
```

### Dashboard Commands
```bash
# From /dashboard directory
npm run dev    # Development server with Turbopack
npm run build  # Production build
npm run lint   # ESLint checking
```

### Go Service Commands
```bash
# From /service directory
go run main.go           # Start service
go mod tidy             # Clean dependencies
go test ./...           # Run tests
```

## Architecture Overview

This is an AI News Auditability Service providing observability for AI-generated content with two main components:

### Backend (Go Service)
- **Framework**: Echo v4 web framework on port 8080
- **Database**: BadgerDB (embedded key-value store) with TTL support
- **API Base**: `/v1` with endpoints for runs, benchmarks, and metrics
- **Configuration**: Viper with environment variables (ANAS_ prefix)

### Frontend (Next.js Dashboard)
- **Framework**: Next.js 15 with App Router and TypeScript
- **Styling**: Tailwind CSS v4
- **Purpose**: Displays run data, summaries, and quality metrics

### Data Flow
1. AI News Processor submits run data via POST `/v1/runs`
2. BadgerDB stores data with configurable TTL (default 7 days)
3. Dashboard fetches latest data from `/v1/runs/latest`
4. Benchmark service evaluates content quality using LLM APIs

## Key File Locations

### Go Service (`/service/`)
- `main.go` - Service entry point
- `internal/api/handlers.go` - API endpoints implementation
- `internal/storage/badger.go` - Database operations
- `internal/models/` - Data structures (PersistedRunData, RunMetadata, etc.)
- `internal/benchmark/benchmark.go` - Benchmark service implementation
- `internal/openai/openai.go` - OpenAI API client for LLM interactions

### Dashboard (`/dashboard/`)
- `src/app/page.tsx` - Main dashboard page
- `src/components/` - React components (MetadataCard, AccordionItem, etc.)
- `src/types/` - TypeScript type definitions

### Benchmark (`/benchmark/`)
- Quality evaluation system for AI-generated content
- Uses OpenAI-compatible APIs for content assessment

## Environment Configuration

Key environment variables:
- `ANAS_RUN_DATA_TTL_HOURS` - Data retention period (default: 168)
- `ANAS_CORS_ALLOWED_ORIGINS` - CORS configuration
- `ANAS_LLM_URL` - LLM API endpoint URL for benchmarking
- `ANAS_LLM_API_KEY` - LLM API key for benchmarking
- `ANAS_LLM_MODEL` - LLM model name (default: gpt-4)

## Database Schema

BadgerDB stores:
- **Run Data**: Complete processing results with metadata
- **Benchmarks**: Quality evaluation results
- **TTL**: Automatic cleanup after configured hours

## API Endpoints

Main endpoints:
- `POST /v1/runs` - Submit new run data
- `GET /v1/runs/latest` - Get most recent run (used by dashboard)
- `GET /v1/runs/{runId}` - Get specific run
- `POST /v1/benchmarks/create/{runId}` - Trigger quality evaluation
- `GET /v1/metrics/*` - Various metrics endpoints

## Testing

The service uses standard Go testing patterns. Run tests with `go test ./...` from the service directory.