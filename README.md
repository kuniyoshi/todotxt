# todo.txt

A Go implementation of the todo.txt format, a simple and powerful task management system using plain text files.

## About

This project implements the todo.txt format specification, which allows you to manage your tasks in a simple, plain text file. The todo.txt format is a set of rules for capturing tasks in a future-proof, human-readable, and tool-friendly way.

## Features

- Simple plain text format for task management
- Priority levels (A-Z)
- Project and context tags (+project @context)
- Due dates and creation dates
- Completion tracking
- List all projects and contexts with task counts
- Archive completed tasks

## Installation

### Using go install

```bash
go install github.com/kuniyoshi/todotxt@latest
```

### Build from source

```bash
git clone https://github.com/kuniyoshi/todotxt
cd todotxt
go build -o todotxt
```

## Usage

### Command Line Interface

```bash
# Add a new task
todotxt add "Buy milk @store"
todotxt add "(A) Call Mom +Family @phone"

# List tasks
todotxt list                  # Show incomplete tasks
todotxt list all              # Show all tasks
todotxt list done             # Show completed tasks
todotxt list +Work            # Filter by project
todotxt list @office          # Filter by context

# Complete a task
todotxt do 1                  # Mark task 1 as complete

# Undo completion
todotxt undo 1                # Mark task 1 as incomplete

# Delete a task
todotxt delete 1              # Remove task 1

# Set priority
todotxt priority 2 B          # Set task 2 to priority B

# Remove priority
todotxt depri 2               # Remove priority from task 2

# List projects and contexts
todotxt projects              # List all projects (incomplete tasks only)
todotxt projects all          # List all projects (including completed)
todotxt contexts              # List all contexts (incomplete tasks only)
todotxt contexts all          # List all contexts (including completed)

# Archive completed tasks
todotxt archive               # Move completed tasks to done.txt

# Help
todotxt help                  # Show usage information
```

### Basic Format

Each line in your todo.txt file represents a single task:

```
(A) Call Mom @phone +Family
(B) Schedule annual checkup +Health @hospital due:2025-02-01
x 2025-01-09 Finish quarterly report +Work @office
2025-01-08 Buy groceries @store
```

### Priority Levels

- `(A)` - Highest priority
- `(B)` - High priority
- `(C)` - Medium priority
- And so on through `(Z)`

### Special Markers

- `x` - Marks a task as complete
- `x 2025-01-09` - Completion date
- `2025-01-08` - Creation date
- `+` - Project tag (e.g., `+Work`)
- `@` - Context tag (e.g., `@office`)
- `key:value` - Custom tags (e.g., `due:2025-01-15`)

### Environment Variables

- `TODO_FILE` - Path to your todo.txt file (default: `~/todo.txt`)
- `DONE_FILE` - Path to your done.txt archive file (default: `~/done.txt`)

Example:
```bash
export TODO_FILE=/path/to/my/tasks.txt
export DONE_FILE=/path/to/my/completed.txt
```

## Development

### Prerequisites

- Go 1.24 or higher

### Building

```bash
go build -o todotxt
```

### Testing

```bash
go test -v        # Run all tests
go test -cover    # Run tests with coverage
```

### Project Structure

```
todotxt/
├── main.go           # CLI entry point
├── todo.go           # Core Todo struct and methods
├── parser.go         # Todo.txt format parser
├── file.go           # File I/O operations
├── commands.go       # CLI command implementations
├── sort.go           # Sorting and filtering functions
├── *_test.go         # Test files
└── README.md         # This file
```

## Todo.txt Format Specification

This implementation follows the todo.txt format rules as described at:
- [Official todo.txt format](https://github.com/todotxt/todo.txt)
- [Qiita Article (Japanese)](https://qiita.com/maedana/items/713390ce590b92fee97f)

## License

[Add your license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.