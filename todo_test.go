package main

import (
	"testing"
	"time"
)

func TestNewTodo(t *testing.T) {
	todo := NewTodo("Test task")

	if todo.Description != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", todo.Description)
	}

	if todo.Complete {
		t.Error("New todo should not be complete")
	}

	if todo.Priority != PriorityNone {
		t.Error("New todo should have no priority")
	}

	if todo.CreationDate == nil {
		t.Error("New todo should have creation date")
	}
}

func TestTodoString(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Todo
		expected string
	}{
		{
			name: "Simple todo",
			setup: func() *Todo {
				return &Todo{
					Description: "Simple task",
				}
			},
			expected: "Simple task",
		},
		{
			name: "Todo with priority",
			setup: func() *Todo {
				return &Todo{
					Priority:    PriorityA,
					Description: "Important task",
				}
			},
			expected: "(A) Important task",
		},
		{
			name: "Completed todo",
			setup: func() *Todo {
				todo := &Todo{
					Complete:    true,
					Description: "Finished task",
				}
				completionDate := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)
				todo.CompletionDate = &completionDate
				return todo
			},
			expected: "x 2025-01-09 Finished task",
		},
		{
			name: "Todo with project and context",
			setup: func() *Todo {
				todo := &Todo{
					Description: "Task with tags",
					Projects:    []string{"Work"},
					Contexts:    []string{"office"},
				}
				return todo
			},
			expected: "Task with tags +Work @office",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := tt.setup()
			result := todo.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestTodoMarkComplete(t *testing.T) {
	todo := &Todo{
		Priority:    PriorityA,
		Description: "Test task",
	}

	todo.MarkComplete()

	if !todo.Complete {
		t.Error("Todo should be marked as complete")
	}

	if todo.Priority != PriorityNone {
		t.Error("Completed todo should have no priority")
	}

	if todo.CompletionDate == nil {
		t.Error("Completed todo should have completion date")
	}
}

func TestTodoMarkUncomplete(t *testing.T) {
	todo := &Todo{
		Complete:    true,
		Description: "Test task",
	}
	completionDate := time.Now()
	todo.CompletionDate = &completionDate

	todo.MarkUncomplete()

	if todo.Complete {
		t.Error("Todo should not be complete")
	}

	if todo.CompletionDate != nil {
		t.Error("Uncompleted todo should have no completion date")
	}
}

func TestAddProject(t *testing.T) {
	todo := &Todo{
		Description: "Test task",
		Projects:    []string{},
	}

	todo.AddProject("Work")
	if len(todo.Projects) != 1 || todo.Projects[0] != "Work" {
		t.Error("Project not added correctly")
	}

	todo.AddProject("Work")
	if len(todo.Projects) != 1 {
		t.Error("Duplicate project should not be added")
	}

	todo.AddProject("Personal")
	if len(todo.Projects) != 2 || todo.Projects[1] != "Personal" {
		t.Error("Second project not added correctly")
	}
}

func TestAddContext(t *testing.T) {
	todo := &Todo{
		Description: "Test task",
		Contexts:    []string{},
	}

	todo.AddContext("home")
	if len(todo.Contexts) != 1 || todo.Contexts[0] != "home" {
		t.Error("Context not added correctly")
	}

	todo.AddContext("home")
	if len(todo.Contexts) != 1 {
		t.Error("Duplicate context should not be added")
	}
}

func TestAddTag(t *testing.T) {
	todo := &Todo{
		Description: "Test task",
		Tags:        make(map[string]string),
	}

	todo.AddTag("due", "2025-01-15")
	if todo.Tags["due"] != "2025-01-15" {
		t.Error("Tag not added correctly")
	}

	todo.AddTag("due", "2025-01-20")
	if todo.Tags["due"] != "2025-01-20" {
		t.Error("Tag should be updated")
	}
}

func TestGetDueDate(t *testing.T) {
	todo := &Todo{
		Description: "Test task",
		Tags:        map[string]string{"due": "2025-01-15"},
	}

	dueDate := todo.GetDueDate()
	if dueDate == nil {
		t.Fatal("Due date should not be nil")
	}

	expected := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	if !dueDate.Equal(expected) {
		t.Errorf("Expected due date %v, got %v", expected, dueDate)
	}

	todo.Tags["due"] = "invalid-date"
	if todo.GetDueDate() != nil {
		t.Error("Invalid date should return nil")
	}
}
