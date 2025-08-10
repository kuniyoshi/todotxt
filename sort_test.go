package main

import (
	"testing"
	"time"
)

func TestSortByPriority(t *testing.T) {
	todos := []*Todo{
		{ID: 1, Priority: PriorityC, Description: "Task C"},
		{ID: 2, Priority: PriorityA, Description: "Task A"},
		{ID: 3, Priority: PriorityNone, Description: "No priority"},
		{ID: 4, Priority: PriorityB, Description: "Task B"},
	}

	SortTodos(todos, SortByPriority)

	if todos[0].Priority != PriorityA {
		t.Error("First task should have priority A")
	}
	if todos[1].Priority != PriorityB {
		t.Error("Second task should have priority B")
	}
	if todos[2].Priority != PriorityC {
		t.Error("Third task should have priority C")
	}
	if todos[3].Priority != PriorityNone {
		t.Error("Last task should have no priority")
	}
}

func TestSortByCreationDate(t *testing.T) {
	date1 := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)

	todos := []*Todo{
		{ID: 1, CreationDate: &date1, Description: "Task 1"},
		{ID: 2, CreationDate: &date2, Description: "Task 2"},
		{ID: 3, CreationDate: nil, Description: "No date"},
		{ID: 4, CreationDate: &date3, Description: "Task 3"},
	}

	SortTodos(todos, SortByCreationDate)

	if !todos[0].CreationDate.Equal(date3) {
		t.Error("First task should have earliest date")
	}
	if !todos[1].CreationDate.Equal(date1) {
		t.Error("Second task should have middle date")
	}
	if !todos[2].CreationDate.Equal(date2) {
		t.Error("Third task should have latest date")
	}
	if todos[3].CreationDate != nil {
		t.Error("Last task should have no date")
	}
}

func TestSortByDueDate(t *testing.T) {
	todos := []*Todo{
		{ID: 1, Description: "Task 1", Tags: map[string]string{"due": "2025-01-20"}},
		{ID: 2, Description: "Task 2", Tags: map[string]string{"due": "2025-01-10"}},
		{ID: 3, Description: "No due", Tags: map[string]string{}},
		{ID: 4, Description: "Task 3", Tags: map[string]string{"due": "2025-01-15"}},
	}

	SortTodos(todos, SortByDueDate)

	if todos[0].Tags["due"] != "2025-01-10" {
		t.Error("First task should have earliest due date")
	}
	if todos[1].Tags["due"] != "2025-01-15" {
		t.Error("Second task should have middle due date")
	}
	if todos[2].Tags["due"] != "2025-01-20" {
		t.Error("Third task should have latest due date")
	}
	if _, hasDue := todos[3].Tags["due"]; hasDue {
		t.Error("Last task should have no due date")
	}
}

func TestFilterOverdue(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	todos := []*Todo{
		{ID: 1, Complete: false, Description: "Overdue", Tags: map[string]string{"due": yesterday}},
		{ID: 2, Complete: false, Description: "Future", Tags: map[string]string{"due": tomorrow}},
		{ID: 3, Complete: true, Description: "Complete overdue", Tags: map[string]string{"due": yesterday}},
		{ID: 4, Complete: false, Description: "No due", Tags: map[string]string{}},
	}

	overdue := FilterOverdue(todos)

	if len(overdue) != 1 {
		t.Errorf("Expected 1 overdue task, got %d", len(overdue))
	}
	if overdue[0].ID != 1 {
		t.Error("Wrong task identified as overdue")
	}
}

func TestFilterToday(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	todos := []*Todo{
		{ID: 1, Complete: false, Description: "Today", Tags: map[string]string{"due": today}},
		{ID: 2, Complete: false, Description: "Tomorrow", Tags: map[string]string{"due": tomorrow}},
		{ID: 3, Complete: false, Description: "Yesterday", Tags: map[string]string{"due": yesterday}},
		{ID: 4, Complete: false, Description: "No due", Tags: map[string]string{}},
	}

	todayTasks := FilterToday(todos)

	if len(todayTasks) != 1 {
		t.Errorf("Expected 1 task for today, got %d", len(todayTasks))
	}
	if todayTasks[0].ID != 1 {
		t.Error("Wrong task identified for today")
	}
}

func TestGroupByProject(t *testing.T) {
	todos := []*Todo{
		{ID: 1, Description: "Task 1", Projects: []string{"Work"}},
		{ID: 2, Description: "Task 2", Projects: []string{"Work", "Urgent"}},
		{ID: 3, Description: "Task 3", Projects: []string{"Personal"}},
		{ID: 4, Description: "Task 4", Projects: []string{}},
	}

	groups := GroupByProject(todos)

	if len(groups["Work"]) != 2 {
		t.Errorf("Expected 2 tasks in Work project, got %d", len(groups["Work"]))
	}
	if len(groups["Personal"]) != 1 {
		t.Errorf("Expected 1 task in Personal project, got %d", len(groups["Personal"]))
	}
	if len(groups["Urgent"]) != 1 {
		t.Errorf("Expected 1 task in Urgent project, got %d", len(groups["Urgent"]))
	}
	if len(groups["No Project"]) != 1 {
		t.Errorf("Expected 1 task with no project, got %d", len(groups["No Project"]))
	}
}

func TestGroupByContext(t *testing.T) {
	todos := []*Todo{
		{ID: 1, Description: "Task 1", Contexts: []string{"home"}},
		{ID: 2, Description: "Task 2", Contexts: []string{"office", "meeting"}},
		{ID: 3, Description: "Task 3", Contexts: []string{"office"}},
		{ID: 4, Description: "Task 4", Contexts: []string{}},
	}

	groups := GroupByContext(todos)

	if len(groups["office"]) != 2 {
		t.Errorf("Expected 2 tasks in office context, got %d", len(groups["office"]))
	}
	if len(groups["home"]) != 1 {
		t.Errorf("Expected 1 task in home context, got %d", len(groups["home"]))
	}
	if len(groups["meeting"]) != 1 {
		t.Errorf("Expected 1 task in meeting context, got %d", len(groups["meeting"]))
	}
	if len(groups["No Context"]) != 1 {
		t.Errorf("Expected 1 task with no context, got %d", len(groups["No Context"]))
	}
}
