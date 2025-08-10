package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTodoFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_todo.txt")

	tf := NewTodoFile(testFile)

	todo1 := NewTodo("First task")
	todo1.Priority = PriorityA
	tf.Add(todo1)

	todo2 := NewTodo("Second task")
	todo2.AddProject("Work")
	tf.Add(todo2)

	todo3 := NewTodo("Third task")
	todo3.AddContext("home")
	tf.Add(todo3)

	if err := tf.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("File should exist after save")
	}

	tf2 := NewTodoFile(testFile)
	if err := tf2.Load(); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if len(tf2.Todos) != 3 {
		t.Errorf("Expected 3 todos, got %d", len(tf2.Todos))
	}

	if tf2.Todos[0].Priority != PriorityA {
		t.Error("First todo should have priority A")
	}

	if len(tf2.Todos[1].Projects) != 1 || tf2.Todos[1].Projects[0] != "Work" {
		t.Error("Second todo should have project 'Work'")
	}
}

func TestTodoFileGetByID(t *testing.T) {
	tf := NewTodoFile("test.txt")

	todo1 := NewTodo("Task 1")
	tf.Add(todo1)

	todo2 := NewTodo("Task 2")
	tf.Add(todo2)

	found := tf.GetByID(1)
	if found == nil {
		t.Fatal("Should find todo with ID 1")
	}
	if found.Description != "Task 1" {
		t.Error("Wrong todo returned")
	}

	notFound := tf.GetByID(99)
	if notFound != nil {
		t.Error("Should return nil for non-existent ID")
	}
}

func TestTodoFileDelete(t *testing.T) {
	tf := NewTodoFile("test.txt")

	todo1 := NewTodo("Task 1")
	tf.Add(todo1)

	todo2 := NewTodo("Task 2")
	tf.Add(todo2)

	todo3 := NewTodo("Task 3")
	tf.Add(todo3)

	if !tf.Delete(2) {
		t.Error("Delete should return true")
	}

	if len(tf.Todos) != 2 {
		t.Errorf("Expected 2 todos after delete, got %d", len(tf.Todos))
	}

	if tf.Todos[0].ID != 1 || tf.Todos[1].ID != 2 {
		t.Error("IDs should be reindexed after delete")
	}

	if tf.Delete(99) {
		t.Error("Delete should return false for non-existent ID")
	}
}

func TestTodoFileSearch(t *testing.T) {
	tf := NewTodoFile("test.txt")

	todo1 := NewTodo("Buy milk")
	tf.Add(todo1)

	todo2 := NewTodo("Call mom")
	todo2.AddProject("Family")
	tf.Add(todo2)

	todo3 := NewTodo("Write report")
	todo3.AddContext("work")
	tf.Add(todo3)

	results := tf.Search("milk")
	if len(results) != 1 || results[0].Description != "Buy milk" {
		t.Error("Should find 'Buy milk'")
	}

	results = tf.Search("family")
	if len(results) != 1 || results[0].ID != 2 {
		t.Error("Should find task with Family project")
	}

	results = tf.Search("work")
	if len(results) != 1 || results[0].ID != 3 {
		t.Error("Should find task with work context")
	}

	results = tf.Search("nonexistent")
	if len(results) != 0 {
		t.Error("Should return empty results for non-matching search")
	}
}

func TestTodoFileFilters(t *testing.T) {
	tf := NewTodoFile("test.txt")

	todo1 := NewTodo("Task 1")
	todo1.AddProject("Work")
	tf.Add(todo1)

	todo2 := NewTodo("Task 2")
	todo2.AddProject("Personal")
	tf.Add(todo2)

	todo3 := NewTodo("Task 3")
	todo3.AddContext("home")
	tf.Add(todo3)

	todo4 := NewTodo("Task 4")
	todo4.Complete = true
	tf.Add(todo4)

	workTodos := tf.FilterByProject("Work")
	if len(workTodos) != 1 || workTodos[0].ID != 1 {
		t.Error("Should filter by project correctly")
	}

	homeTodos := tf.FilterByContext("home")
	if len(homeTodos) != 1 || homeTodos[0].ID != 3 {
		t.Error("Should filter by context correctly")
	}

	completed := tf.GetCompleted()
	if len(completed) != 1 || completed[0].ID != 4 {
		t.Error("Should get completed todos")
	}

	incomplete := tf.GetIncomplete()
	if len(incomplete) != 3 {
		t.Error("Should get incomplete todos")
	}
}
