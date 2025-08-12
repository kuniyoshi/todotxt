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
	todo, _ := ParseTodo(description)

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
	if len(args) > 0 {
		return helpForCommand(args[0])
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    TODO.TXT TASK MANAGER                       ║")
	fmt.Println("║              Simple and powerful task management               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  todotxt [command] [arguments]")
	fmt.Println("  todotxt --help               Show this help message")
	fmt.Println("  todotxt help <command>       Show help for a specific command")
	fmt.Println()
	fmt.Println("TASK MANAGEMENT:")
	fmt.Println("  add <task>               Add a new task")
	fmt.Println("  list, ls [filter]        List tasks (default: incomplete)")
	fmt.Println("  do, done <ID>            Mark task as complete")
	fmt.Println("  undo <ID>                Mark task as incomplete")
	fmt.Println("  delete, rm <ID>          Delete a task")
	fmt.Println()
	fmt.Println("PRIORITY MANAGEMENT:")
	fmt.Println("  priority, pri <ID> <A-Z> Set task priority (A=highest)")
	fmt.Println("  depri <ID>               Remove task priority")
	fmt.Println()
	fmt.Println("ORGANIZATION:")
	fmt.Println("  projects, proj [all]     List all projects with task counts")
	fmt.Println("  contexts, ctx [all]      List all contexts with task counts")
	fmt.Println("  archive                  Move completed tasks to done.txt")
	fmt.Println()
	fmt.Println("LIST FILTERS:")
	fmt.Println("  list                     Show incomplete tasks")
	fmt.Println("  list all                 Show all tasks")
	fmt.Println("  list done                Show completed tasks only")
	fmt.Println("  list +Project            Filter by project")
	fmt.Println("  list @Context            Filter by context")
	fmt.Println("  list <search>            Search in task descriptions")
	fmt.Println()
	fmt.Println("TASK FORMAT:")
	fmt.Println("  (A) Task description +project @context key:value")
	fmt.Println()
	fmt.Println("  Special markers:")
	fmt.Println("    (A-Z)        Priority level")
	fmt.Println("    x            Completed task")
	fmt.Println("    +project     Project tag")
	fmt.Println("    @context     Context tag")
	fmt.Println("    due:date     Due date (format: YYYY-MM-DD)")
	fmt.Println("    key:value    Custom metadata")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  todotxt add \"(A) Call Mom +Family @phone\"")
	fmt.Println("  todotxt add \"Submit report +Work @office due:2025-01-15\"")
	fmt.Println("  todotxt list +Work")
	fmt.Println("  todotxt do 3")
	fmt.Println("  todotxt priority 5 B")
	fmt.Println()
	fmt.Println("ENVIRONMENT:")
	fmt.Println("  TODO_FILE     Path to todo.txt (default: ~/todo.txt)")
	fmt.Println("  DONE_FILE     Path to done.txt (default: ~/done.txt)")
	fmt.Println()
	fmt.Println("For more information on a specific command, run:")
	fmt.Println("  todotxt help <command>")

	return nil
}

func helpForCommand(cmd string) error {
	helps := map[string]string{
		"add": `ADD COMMAND - Add a new task

USAGE:
  todotxt add <task description>

DESCRIPTION:
  Adds a new task to your todo.txt file. The task can include priority,
  projects, contexts, and custom tags.

EXAMPLES:
  todotxt add "Buy milk"
  todotxt add "(A) Important meeting +Work @office"
  todotxt add "Submit report +Work due:2025-01-15"
  todotxt add "(B) Call dentist @phone +Health"

TASK FORMAT:
  - (A-Z)      Priority (A is highest)
  - +project   Project tag (can have multiple)
  - @context   Context tag (can have multiple)
  - key:value  Custom tags (e.g., due:2025-01-15)`,

		"list": `LIST COMMAND - Display tasks

USAGE:
  todotxt list [filter]
  todotxt ls [filter]

DESCRIPTION:
  Lists tasks from your todo.txt file. By default shows incomplete tasks.

FILTERS:
  (none)       Show incomplete tasks
  all          Show all tasks
  done         Show completed tasks only
  +Project     Filter by project
  @Context     Filter by context
  <search>     Search in descriptions

EXAMPLES:
  todotxt list                 # Show incomplete tasks
  todotxt list all            # Show all tasks
  todotxt list done           # Show completed tasks
  todotxt list +Work          # Show tasks in Work project
  todotxt list @home          # Show tasks in home context
  todotxt list "report"       # Search for "report"`,

		"do": `DO/DONE COMMAND - Mark task as complete

USAGE:
  todotxt do <ID>
  todotxt done <ID>
  todotxt complete <ID>

DESCRIPTION:
  Marks a task as complete. This adds an 'x' marker and completion date
  to the task, and removes any priority.

EXAMPLES:
  todotxt do 3
  todotxt done 1
  todotxt complete 5`,

		"priority": `PRIORITY COMMAND - Set task priority

USAGE:
  todotxt priority <ID> <A-Z>
  todotxt pri <ID> <A-Z>

DESCRIPTION:
  Sets or changes the priority of a task. Priority ranges from A (highest)
  to Z (lowest). Completed tasks cannot have priorities.

EXAMPLES:
  todotxt priority 3 A        # Set highest priority
  todotxt pri 5 C            # Set medium priority
  todotxt priority 2 Z       # Set lowest priority`,

		"projects": `PROJECTS COMMAND - List all projects

USAGE:
  todotxt projects [all]
  todotxt proj [all]

DESCRIPTION:
  Lists all unique projects found in tasks, along with the count of
  tasks in each project. By default shows projects from incomplete
  tasks only.

OPTIONS:
  all          Include completed tasks in counts

EXAMPLES:
  todotxt projects            # Projects from incomplete tasks
  todotxt projects all        # Projects from all tasks`,

		"contexts": `CONTEXTS COMMAND - List all contexts

USAGE:
  todotxt contexts [all]
  todotxt ctx [all]

DESCRIPTION:
  Lists all unique contexts found in tasks, along with the count of
  tasks in each context. By default shows contexts from incomplete
  tasks only.

OPTIONS:
  all          Include completed tasks in counts

EXAMPLES:
  todotxt contexts            # Contexts from incomplete tasks
  todotxt contexts all        # Contexts from all tasks`,

		"archive": `ARCHIVE COMMAND - Archive completed tasks

USAGE:
  todotxt archive

DESCRIPTION:
  Moves all completed tasks from todo.txt to done.txt. This helps keep
  your active todo list clean and focused on current tasks.

NOTES:
  - Completed tasks are appended to done.txt
  - Original completion dates are preserved
  - Tasks are removed from todo.txt after archiving

EXAMPLE:
  todotxt archive`,

		"delete": `DELETE COMMAND - Remove a task

USAGE:
  todotxt delete <ID>
  todotxt rm <ID>
  todotxt del <ID>

DESCRIPTION:
  Permanently removes a task from your todo.txt file. This action
  cannot be undone.

EXAMPLES:
  todotxt delete 3
  todotxt rm 5
  todotxt del 1`,

		"undo": `UNDO COMMAND - Mark task as incomplete

USAGE:
  todotxt undo <ID>
  todotxt undone <ID>

DESCRIPTION:
  Marks a completed task as incomplete again. This removes the 'x'
  marker and completion date from the task.

EXAMPLES:
  todotxt undo 3
  todotxt undone 5`,

		"depri": `DEPRI COMMAND - Remove task priority

USAGE:
  todotxt depri <ID>

DESCRIPTION:
  Removes the priority from a task. The task will no longer have
  a priority marker (A-Z).

EXAMPLES:
  todotxt depri 3
  todotxt depri 1`,
	}

	// Check for command aliases
	aliases := map[string]string{
		"ls":       "list",
		"done":     "do",
		"complete": "do",
		"rm":       "delete",
		"del":      "delete",
		"pri":      "priority",
		"proj":     "projects",
		"ctx":      "contexts",
	}

	if alias, ok := aliases[cmd]; ok {
		cmd = alias
	}

	if help, ok := helps[cmd]; ok {
		fmt.Println(help)
		return nil
	}

	return fmt.Errorf("no help available for command: %s", cmd)
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
	// Check for --help or -h before parsing flags
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" || arg == "-help" {
			return "help", []string{}
		}
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return "list", []string{}
	}

	return args[0], args[1:]
}
