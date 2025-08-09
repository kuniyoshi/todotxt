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

## Installation

```bash
go get todotxt.kuniyoshi.github.com
```

## Usage

### Basic Format

Each line in your todo.txt file represents a single task:

```
(A) Call Mom @phone +Family
(B) Schedule annual checkup +Health @hospital due:2025-02-01
x 2025-01-09 (C) Finish quarterly report +Work @office
```

### Priority Levels

- `(A)` - Highest priority
- `(B)` - High priority
- `(C)` - Medium priority
- And so on...

### Special Markers

- `x` - Marks a task as complete
- `+` - Project tag (e.g., `+Work`)
- `@` - Context tag (e.g., `@office`)
- `due:` - Due date (e.g., `due:2025-01-15`)

## Development

### Prerequisites

- Go 1.24 or higher

### Building

```bash
go build
```

### Running

```bash
./todotxt.kuniyoshi.github.com
```

## Todo.txt Format Specification

This implementation follows the todo.txt format rules as described at:
- [Official todo.txt format](https://github.com/todotxt/todo.txt)
- [Qiita Article (Japanese)](https://qiita.com/maedana/items/713390ce590b92fee97f)

## License

[Add your license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.