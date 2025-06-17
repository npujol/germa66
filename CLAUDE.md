# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

germa66 is a Go CLI application for converting and importing dictionary files to MeiliSearch. The application converts BGL dictionary files to CSV format using pyglossary, then imports them into a MeiliSearch index in configurable batches.

## Architecture

The application follows a clean architecture pattern with these key layers:

- **Main Entry Point** (`main.go`): Cobra CLI setup with configuration and input file flags
- **Config Layer** (`internal/config/`): Environment-based configuration management using Viper
- **MeiliClient Layer** (`internal/meiliclient/`): MeiliSearch operations and batch processing logic  
- **Models Layer** (`internal/models/`): Card data structure representing dictionary entries
- **Utils Layer** (`internal/utils/`): File operations, logging, and pyglossary integration

## Key Dependencies

- **External Tools**: Requires `pyglossary` command-line tool for BGL to CSV conversion
- **MeiliSearch**: Uses meilisearch-go client for search engine operations
- **CLI Framework**: Built with Cobra for command-line interface
- **Configuration**: Viper for environment variable management
- **Logging**: Logrus for structured logging

## Development Commands

### Building and Running
```bash
# Build the application
go build -o main -buildvcs=false

# Run with configuration
./main --config ./path/to/config.env --input ./import/deutsch_spanisch.BGL

# Run tests
go test ./...

# Run specific test package
go test ./internal/config
go test ./internal/meiliclient
```

### Docker Development
```bash
# Build Docker image
docker build -t germa66 .

# Run with Docker Compose
docker-compose up

# Run testing environment
docker-compose -f docker-compose.testing.yml up
```

### Configuration Requirements

The application requires these environment variables:
- `MEILISEARCH_HOST`: MeiliSearch server URL
- `MEILISEARCH_API_KEY`: MeiliSearch API key
- `DEBUG`: Enable debug logging (optional, boolean)

Default batch size is 20000 documents per upload batch.

## Key Implementation Details

### Dictionary Import Flow
1. CLI accepts BGL input file and configuration path
2. Configuration validation ensures required MeiliSearch credentials
3. pyglossary converts BGL to CSV format
4. CSV parsing with batch processing (configurable batch size)
5. MeiliSearch document upload with task tracking
6. Comprehensive logging throughout the process

### Data Model
Cards represent dictionary entries with fields:
- `id`: Unique identifier (uses word as ID)
- `word`: Dictionary term
- `description`: Definition/translation
- `backend`: Source dictionary file name

### Error Handling
- Configuration validation with specific error types
- CSV parsing with row-level error handling and skipping
- MeiliSearch operation error handling with task UIDs
- Comprehensive logging at debug, info, warn, and error levels

## Testing

The project uses Go's standard testing framework with testify for assertions. Test files follow the `*_test.go` naming convention and are located alongside their corresponding source files.