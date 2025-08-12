# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Essential Commands
```bash
# Build the binary
make build              # Creates bin/todotxt
go build -o todotxt     # Alternative direct build

# Run all quality checks before committing
make check              # Runs fmt, vet, staticcheck, lint, and tests

# Testing
make test               # Run all tests
make test-coverage      # Generate coverage.html report
go test -v ./...        # Run tests with verbose output
go test -run TestAdd    # Run specific test

# Code quality
make fmt                # Format all Go code
make staticcheck        # Run static analysis
make lint               # Run linter (requires golangci-lint)

# Cross-platform builds
make build-all          # Build for Linux, Windows, macOS (Intel & ARM)
make release            # Create distribution archives
```

## Architecture Overview

### Command Flow Architecture
The codebase follows a layered architecture where commands flow through distinct layers:
1. **main.go** → Entry point that initializes global TodoFile and routes commands
2. **commands.go** → Command layer that handles all CLI operations via executeCommand()
3. **todo.go + file.go** → Domain layer with Todo struct and TodoFile operations
4. **parser.go** → Parsing layer that converts text ↔ Todo objects

### Global State Pattern
The application uses a global `todoFile *TodoFile` variable initialized once at startup. All commands operate on this shared instance, which:
- Loads the entire todo.txt file into memory at startup
- Maintains an in-memory array of Todo pointers
- Saves the entire file after each modification
- Auto-reindexes IDs after deletions

### Command Routing System
Commands are dispatched through a map-based router with extensive aliasing:
- Primary commands are keys in the executeCommand switch
- Aliases map to primary commands (e.g., "ls"→"list", "do"→"done"→"complete")
- Each command validates arguments then operates on the global TodoFile
- Commands return errors that propagate to main for user-friendly display

### Parser Architecture
The parser uses a stateful regex-based approach:
1. Pre-compiled regex patterns for each todo.txt element
2. Sequential parsing: completion → priority → dates → description → tags
3. Preserves original raw text alongside parsed data
4. Projects (+) and contexts (@) are extracted into separate arrays
5. Generic key:value tags go into a map for extensibility

### ID Management System
IDs are not stored but computed dynamically:
- IDs are 1-based sequential numbers for user display
- After deletion, all todos are reindexed to maintain sequential IDs
- GetByID() searches linearly through the array
- This ensures consistent, gap-free numbering for CLI usage

### File I/O Strategy
All file operations follow a load-modify-save pattern:
- TODO_FILE and DONE_FILE paths from environment or defaults
- Load() reads entire file and parses all lines at once
- Save() writes all todos back (full file rewrite)
- Archive moves completed tasks to done.txt
- No incremental updates - simpler but requires full rewrites

## Key Implementation Details

### Date Handling
- Dates use pointers (*time.Time) to distinguish unset from zero values
- Format: YYYY-MM-DD (todo.txt standard)
- Completion automatically adds today's date
- Creation date preserved from original text

### Priority System
- Custom Priority type with constants PriorityNone through PriorityZ
- Completed tasks have priority removed
- Sort order: A (highest) → Z (lowest) → None

### Sorting and Filtering
- SortTodos() provides multi-criteria sorting
- Filter functions for overdue, today, this week
- Projects/contexts extracted and counted across all todos
- Search is case-insensitive substring matching

## Testing Approach

Tests use table-driven patterns with temporary files:
```go
// Most tests follow this pattern
testCases := []struct {
    name     string
    input    string
    expected *Todo
}{...}

// File tests use temp directories
tmpDir := t.TempDir()
todoPath := filepath.Join(tmpDir, "test_todo.txt")
```

Run specific tests during development:
```bash
go test -v -run TestParseTodo     # Test parser
go test -v -run TestTodoFile      # Test file operations
go test -v -run TestSort          # Test sorting
```

## Module Path
The module is `github.com/kuniyoshi/todotxt` - ensure this matches in go.mod when making changes to module dependencies or structure.