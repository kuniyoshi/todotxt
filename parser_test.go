package main

import (
	"testing"
	"time"
)

func TestParseTodo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(*testing.T, *Todo)
	}{
		{
			name:  "Simple task",
			input: "Buy milk",
			validate: func(t *testing.T, todo *Todo) {
				if todo.Description != "Buy milk" {
					t.Errorf("Expected description 'Buy milk', got '%s'", todo.Description)
				}
				if todo.Complete {
					t.Error("Should not be complete")
				}
			},
		},
		{
			name:  "Task with priority",
			input: "(A) Call Mom",
			validate: func(t *testing.T, todo *Todo) {
				if todo.Priority != PriorityA {
					t.Errorf("Expected priority A, got %c", todo.Priority)
				}
				if todo.Description != "Call Mom" {
					t.Errorf("Expected description 'Call Mom', got '%s'", todo.Description)
				}
			},
		},
		{
			name:  "Completed task",
			input: "x 2025-01-09 2025-01-08 Write tests",
			validate: func(t *testing.T, todo *Todo) {
				if !todo.Complete {
					t.Error("Should be complete")
				}
				if todo.CompletionDate == nil {
					t.Fatal("Should have completion date")
				}
				expected := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)
				if !todo.CompletionDate.Equal(expected) {
					t.Errorf("Expected completion date %v, got %v", expected, todo.CompletionDate)
				}
				if todo.CreationDate == nil {
					t.Fatal("Should have creation date")
				}
				expected = time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC)
				if !todo.CreationDate.Equal(expected) {
					t.Errorf("Expected creation date %v, got %v", expected, todo.CreationDate)
				}
			},
		},
		{
			name:  "Task with project and context",
			input: "Review PR +Work @office",
			validate: func(t *testing.T, todo *Todo) {
				if todo.Description != "Review PR" {
					t.Errorf("Expected description 'Review PR', got '%s'", todo.Description)
				}
				if len(todo.Projects) != 1 || todo.Projects[0] != "Work" {
					t.Error("Should have project 'Work'")
				}
				if len(todo.Contexts) != 1 || todo.Contexts[0] != "office" {
					t.Error("Should have context 'office'")
				}
			},
		},
		{
			name:  "Task with due date",
			input: "Submit report due:2025-01-15",
			validate: func(t *testing.T, todo *Todo) {
				if todo.Description != "Submit report" {
					t.Errorf("Expected description 'Submit report', got '%s'", todo.Description)
				}
				if todo.Tags["due"] != "2025-01-15" {
					t.Error("Should have due date tag")
				}
			},
		},
		{
			name:  "Task with multiple projects and contexts",
			input: "Meeting about project +Work +Planning @office @conference-room",
			validate: func(t *testing.T, todo *Todo) {
				if len(todo.Projects) != 2 {
					t.Errorf("Expected 2 projects, got %d", len(todo.Projects))
				}
				if len(todo.Contexts) != 2 {
					t.Errorf("Expected 2 contexts, got %d", len(todo.Contexts))
				}
			},
		},
		{
			name:  "Task with priority and creation date",
			input: "(B) 2025-01-08 Schedule dentist appointment",
			validate: func(t *testing.T, todo *Todo) {
				if todo.Priority != PriorityB {
					t.Errorf("Expected priority B, got %c", todo.Priority)
				}
				if todo.CreationDate == nil {
					t.Fatal("Should have creation date")
				}
				expected := time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC)
				if !todo.CreationDate.Equal(expected) {
					t.Errorf("Expected creation date %v, got %v", expected, todo.CreationDate)
				}
			},
		},
		{
			name:  "Empty line",
			input: "",
			validate: func(t *testing.T, todo *Todo) {
				if todo != nil {
					t.Error("Empty line should return nil")
				}
			},
		},
		{
			name:  "Whitespace only",
			input: "   ",
			validate: func(t *testing.T, todo *Todo) {
				if todo != nil {
					t.Error("Whitespace line should return nil")
				}
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := ParseTodo(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			tt.validate(t, todo)
		})
	}
}

func TestParseTodos(t *testing.T) {
	lines := []string{
		"(A) Call Mom",
		"",
		"Buy milk @store",
		"x 2025-01-09 Finish report +Work",
		"   ",
	}
	
	todos, err := ParseTodos(lines)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(todos))
	}
	
	if todos[0].Priority != PriorityA {
		t.Error("First todo should have priority A")
	}
	
	if todos[0].ID != 1 {
		t.Error("First todo should have ID 1")
	}
	
	if len(todos[1].Contexts) != 1 || todos[1].Contexts[0] != "store" {
		t.Error("Second todo should have context 'store'")
	}
	
	if !todos[2].Complete {
		t.Error("Third todo should be complete")
	}
}