package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Execute     func([]string) error
}

var todoFile *TodoFile

func initTodoFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	todoPath := os.Getenv("TODO_FILE")
	if todoPath == "" {
		todoPath = homeDir + "/todo.txt"
	}

	todoFile = NewTodoFile(todoPath)
	if err := todoFile.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading todo file: %v\n", err)
		os.Exit(1)
	}
}

func saveFile() error {
	return todoFile.Save()
}

func addCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task description provided")
	}

	description := strings.Join(args, " ")
	todo := NewTodo(description)

	todo, _ = ParseTodo(description)

	todoFile.Add(todo)

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Added: %s\n", todo.String())
	return nil
}

func listCommand(args []string) error {
	todos := todoFile.GetIncomplete()

	if len(args) > 0 {
		switch args[0] {
		case "all":
			todos = todoFile.Todos
		case "done":
			todos = todoFile.GetCompleted()
		default:
			if strings.HasPrefix(args[0], "+") {
				project := strings.TrimPrefix(args[0], "+")
				todos = todoFile.FilterByProject(project)
			} else if strings.HasPrefix(args[0], "@") {
				context := strings.TrimPrefix(args[0], "@")
				todos = todoFile.FilterByContext(context)
			} else {
				query := strings.Join(args, " ")
				todos = todoFile.Search(query)
			}
		}
	}

	if len(todos) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for _, todo := range todos {
		status := " "
		if todo.Complete {
			status = "x"
		}
		fmt.Printf("[%s] %3d: %s\n", status, todo.ID, todo.String())
	}

	return nil
}

func completeCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task ID provided")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	todo := todoFile.GetByID(id)
	if todo == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	todo.MarkComplete()

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Completed: %s\n", todo.String())
	return nil
}

func uncompleteCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task ID provided")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	todo := todoFile.GetByID(id)
	if todo == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	todo.MarkUncomplete()

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Uncompleted: %s\n", todo.String())
	return nil
}

func deleteCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task ID provided")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	todo := todoFile.GetByID(id)
	if todo == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	description := todo.String()

	if !todoFile.Delete(id) {
		return fmt.Errorf("failed to delete task with ID %d", id)
	}

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Deleted: %s\n", description)
	return nil
}

func priorityCommand(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: priority <ID> <A-Z>")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	todo := todoFile.GetByID(id)
	if todo == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	priorityStr := strings.ToUpper(args[1])
	if len(priorityStr) != 1 || priorityStr[0] < 'A' || priorityStr[0] > 'Z' {
		return fmt.Errorf("invalid priority: %s (must be A-Z)", args[1])
	}

	todo.SetPriority(Priority(priorityStr[0]))

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Updated priority: %s\n", todo.String())
	return nil
}

func depriCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no task ID provided")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}

	todo := todoFile.GetByID(id)
	if todo == nil {
		return fmt.Errorf("task with ID %d not found", id)
	}

	todo.SetPriority(PriorityNone)

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Removed priority: %s\n", todo.String())
	return nil
}

func archiveCommand(args []string) error {
	homeDir, _ := os.UserHomeDir()
	archivePath := os.Getenv("DONE_FILE")
	if archivePath == "" {
		archivePath = homeDir + "/done.txt"
	}

	archiveFile := NewTodoFile(archivePath)
	if err := archiveFile.Load(); err != nil {
		return fmt.Errorf("failed to load archive file: %w", err)
	}

	completed := todoFile.GetCompleted()
	for _, todo := range completed {
		archiveFile.Add(todo)
	}

	if err := archiveFile.Save(); err != nil {
		return fmt.Errorf("failed to save archive file: %w", err)
	}

	var remaining []*Todo
	for _, todo := range todoFile.Todos {
		if !todo.Complete {
			remaining = append(remaining, todo)
		}
	}
	todoFile.Todos = remaining
	todoFile.reindexTodos()

	if err := saveFile(); err != nil {
		return fmt.Errorf("failed to save todo file: %w", err)
	}

	fmt.Printf("Archived %d completed tasks to %s\n", len(completed), archivePath)
	return nil
}

func projectsCommand(args []string) error {
	projectMap := make(map[string]int)

	todos := todoFile.Todos
	if len(args) > 0 && args[0] == "all" {
		// Include completed tasks
	} else {
		// Only count incomplete tasks by default
		todos = todoFile.GetIncomplete()
	}

	for _, todo := range todos {
		if len(todo.Projects) == 0 {
			projectMap["(none)"]++
		} else {
			for _, project := range todo.Projects {
				projectMap[project]++
			}
		}
	}

	if len(projectMap) == 0 {
		fmt.Println("No projects found.")
		return nil
	}

	// Sort projects alphabetically
	var projects []string
	for project := range projectMap {
		projects = append(projects, project)
	}
	sort.Strings(projects)

	fmt.Println("Projects:")
	for _, project := range projects {
		count := projectMap[project]
		if project == "(none)" {
			fmt.Printf("  %s: %d task(s)\n", project, count)
		} else {
			fmt.Printf("  +%s: %d task(s)\n", project, count)
		}
	}

	return nil
}

func contextsCommand(args []string) error {
	contextMap := make(map[string]int)

	todos := todoFile.Todos
	if len(args) > 0 && args[0] == "all" {
		// Include completed tasks
	} else {
		// Only count incomplete tasks by default
		todos = todoFile.GetIncomplete()
	}

	for _, todo := range todos {
		if len(todo.Contexts) == 0 {
			contextMap["(none)"]++
		} else {
			for _, context := range todo.Contexts {
				contextMap[context]++
			}
		}
	}

	if len(contextMap) == 0 {
		fmt.Println("No contexts found.")
		return nil
	}

	// Sort contexts alphabetically
	var contexts []string
	for context := range contextMap {
		contexts = append(contexts, context)
	}
	sort.Strings(contexts)

	fmt.Println("Contexts:")
	for _, context := range contexts {
		count := contextMap[context]
		if context == "(none)" {
			fmt.Printf("  %s: %d task(s)\n", context, count)
		} else {
			fmt.Printf("  @%s: %d task(s)\n", context, count)
		}
	}

	return nil
}

func helpCommand(args []string) error {
	fmt.Println("todo.txt - Simple and powerful task management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  todo <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <task>           Add a new task")
	fmt.Println("  list [filter]        List tasks (all, done, +project, @context, or search)")
	fmt.Println("  do <ID>              Mark task as complete")
	fmt.Println("  undo <ID>            Mark task as incomplete")
	fmt.Println("  delete <ID>          Delete a task")
	fmt.Println("  priority <ID> <A-Z>  Set task priority")
	fmt.Println("  depri <ID>           Remove task priority")
	fmt.Println("  projects [all]       List all projects")
	fmt.Println("  contexts [all]       List all contexts")
	fmt.Println("  archive              Move completed tasks to done.txt")
	fmt.Println("  help                 Show this help message")
	fmt.Println()
	fmt.Println("Task Format:")
	fmt.Println("  (A) Task description +project @context due:2025-01-15")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  TODO_FILE            Path to todo.txt file (default: ~/todo.txt)")
	fmt.Println("  DONE_FILE            Path to done.txt file (default: ~/done.txt)")

	return nil
}

func executeCommand(name string, args []string) error {
	commands := map[string]func([]string) error{
		"add":      addCommand,
		"list":     listCommand,
		"ls":       listCommand,
		"do":       completeCommand,
		"done":     completeCommand,
		"complete": completeCommand,
		"undo":     uncompleteCommand,
		"undone":   uncompleteCommand,
		"delete":   deleteCommand,
		"del":      deleteCommand,
		"rm":       deleteCommand,
		"priority": priorityCommand,
		"pri":      priorityCommand,
		"depri":    depriCommand,
		"projects": projectsCommand,
		"proj":     projectsCommand,
		"contexts": contextsCommand,
		"ctx":      contextsCommand,
		"archive":  archiveCommand,
		"help":     helpCommand,
	}

	if cmd, ok := commands[name]; ok {
		return cmd(args)
	}

	return fmt.Errorf("unknown command: %s", name)
}

func parseArgs() (string, []string) {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return "list", []string{}
	}

	return args[0], args[1:]
}
