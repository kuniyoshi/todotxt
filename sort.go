package main

import (
	"sort"
	"time"
)

type SortBy int

const (
	SortByID SortBy = iota
	SortByPriority
	SortByCreationDate
	SortByDueDate
	SortByDescription
	SortByComplete
)

func SortTodos(todos []*Todo, sortBy SortBy) {
	switch sortBy {
	case SortByPriority:
		sort.Slice(todos, func(i, j int) bool {
			if todos[i].Priority == PriorityNone && todos[j].Priority == PriorityNone {
				return todos[i].ID < todos[j].ID
			}
			if todos[i].Priority == PriorityNone {
				return false
			}
			if todos[j].Priority == PriorityNone {
				return true
			}
			if todos[i].Priority == todos[j].Priority {
				return todos[i].ID < todos[j].ID
			}
			return todos[i].Priority < todos[j].Priority
		})
	case SortByCreationDate:
		sort.Slice(todos, func(i, j int) bool {
			if todos[i].CreationDate == nil && todos[j].CreationDate == nil {
				return todos[i].ID < todos[j].ID
			}
			if todos[i].CreationDate == nil {
				return false
			}
			if todos[j].CreationDate == nil {
				return true
			}
			if todos[i].CreationDate.Equal(*todos[j].CreationDate) {
				return todos[i].ID < todos[j].ID
			}
			return todos[i].CreationDate.Before(*todos[j].CreationDate)
		})
	case SortByDueDate:
		sort.Slice(todos, func(i, j int) bool {
			iDue := todos[i].GetDueDate()
			jDue := todos[j].GetDueDate()
			
			if iDue == nil && jDue == nil {
				return todos[i].ID < todos[j].ID
			}
			if iDue == nil {
				return false
			}
			if jDue == nil {
				return true
			}
			if iDue.Equal(*jDue) {
				return todos[i].ID < todos[j].ID
			}
			return iDue.Before(*jDue)
		})
	case SortByDescription:
		sort.Slice(todos, func(i, j int) bool {
			if todos[i].Description == todos[j].Description {
				return todos[i].ID < todos[j].ID
			}
			return todos[i].Description < todos[j].Description
		})
	case SortByComplete:
		sort.Slice(todos, func(i, j int) bool {
			if todos[i].Complete == todos[j].Complete {
				return todos[i].ID < todos[j].ID
			}
			return !todos[i].Complete
		})
	default:
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].ID < todos[j].ID
		})
	}
}

func FilterOverdue(todos []*Todo) []*Todo {
	var results []*Todo
	now := time.Now()
	
	for _, todo := range todos {
		if !todo.Complete {
			if due := todo.GetDueDate(); due != nil && due.Before(now) {
				results = append(results, todo)
			}
		}
	}
	
	return results
}

func FilterToday(todos []*Todo) []*Todo {
	var results []*Todo
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	
	for _, todo := range todos {
		if !todo.Complete {
			if due := todo.GetDueDate(); due != nil {
				if !due.Before(today) && due.Before(tomorrow) {
					results = append(results, todo)
				}
			}
		}
	}
	
	return results
}

func FilterThisWeek(todos []*Todo) []*Todo {
	var results []*Todo
	now := time.Now()
	weekFromNow := now.AddDate(0, 0, 7)
	
	for _, todo := range todos {
		if !todo.Complete {
			if due := todo.GetDueDate(); due != nil {
				if !due.Before(now) && due.Before(weekFromNow) {
					results = append(results, todo)
				}
			}
		}
	}
	
	return results
}

func GroupByProject(todos []*Todo) map[string][]*Todo {
	groups := make(map[string][]*Todo)
	groups["No Project"] = []*Todo{}
	
	for _, todo := range todos {
		if len(todo.Projects) == 0 {
			groups["No Project"] = append(groups["No Project"], todo)
		} else {
			for _, project := range todo.Projects {
				groups[project] = append(groups[project], todo)
			}
		}
	}
	
	if len(groups["No Project"]) == 0 {
		delete(groups, "No Project")
	}
	
	return groups
}

func GroupByContext(todos []*Todo) map[string][]*Todo {
	groups := make(map[string][]*Todo)
	groups["No Context"] = []*Todo{}
	
	for _, todo := range todos {
		if len(todo.Contexts) == 0 {
			groups["No Context"] = append(groups["No Context"], todo)
		} else {
			for _, context := range todo.Contexts {
				groups[context] = append(groups[context], todo)
			}
		}
	}
	
	if len(groups["No Context"]) == 0 {
		delete(groups, "No Context")
	}
	
	return groups
}