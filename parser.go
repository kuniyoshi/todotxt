package main

import (
	"regexp"
	"strings"
	"time"
)

var (
	completeRegex     = regexp.MustCompile(`^x `)
	priorityRegex     = regexp.MustCompile(`^\(([A-Z])\) `)
	dateRegex         = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	projectRegex      = regexp.MustCompile(`\+(\S+)`)
	contextRegex      = regexp.MustCompile(`@(\S+)`)
	tagRegex          = regexp.MustCompile(`(\w+):([^\s]+)`)
)

func ParseTodo(line string) (*Todo, error) {
	if strings.TrimSpace(line) == "" {
		return nil, nil
	}
	
	todo := &Todo{
		Projects: []string{},
		Contexts: []string{},
		Tags:     make(map[string]string),
		Raw:      line,
	}
	
	remaining := line
	
	if completeRegex.MatchString(remaining) {
		todo.Complete = true
		remaining = completeRegex.ReplaceAllString(remaining, "")
		
		dates := dateRegex.FindAllString(remaining, 2)
		if len(dates) > 0 {
			if completionDate, err := time.Parse("2006-01-02", dates[0]); err == nil {
				todo.CompletionDate = &completionDate
				remaining = strings.Replace(remaining, dates[0], "", 1)
				remaining = strings.TrimSpace(remaining)
			}
			
			if len(dates) > 1 {
				if creationDate, err := time.Parse("2006-01-02", dates[1]); err == nil {
					todo.CreationDate = &creationDate
					remaining = strings.Replace(remaining, dates[1], "", 1)
					remaining = strings.TrimSpace(remaining)
				}
			}
		}
	}
	
	if !todo.Complete {
		if match := priorityRegex.FindStringSubmatch(remaining); len(match) > 1 {
			todo.Priority = Priority(match[1][0])
			remaining = priorityRegex.ReplaceAllString(remaining, "")
		}
		
		// Only look for creation date at the beginning of remaining text
		// to avoid matching dates that are part of tags
		words := strings.Fields(remaining)
		if len(words) > 0 {
			if creationDate, err := time.Parse("2006-01-02", words[0]); err == nil {
				todo.CreationDate = &creationDate
				remaining = strings.Replace(remaining, words[0], "", 1)
				remaining = strings.TrimSpace(remaining)
			}
		}
	}
	
	for _, match := range projectRegex.FindAllStringSubmatch(line, -1) {
		if len(match) > 1 {
			todo.AddProject(match[1])
		}
	}
	
	for _, match := range contextRegex.FindAllStringSubmatch(line, -1) {
		if len(match) > 1 {
			todo.AddContext(match[1])
		}
	}
	
	for _, match := range tagRegex.FindAllStringSubmatch(line, -1) {
		if len(match) > 2 {
			key := match[1]
			value := match[2]
			if key != "" && value != "" {
				todo.AddTag(key, value)
			}
		}
	}
	
	description := remaining
	
	for _, match := range projectRegex.FindAllStringSubmatch(remaining, -1) {
		if len(match) > 1 {
			description = strings.Replace(description, match[0], "", -1)
		}
	}
	
	for _, match := range contextRegex.FindAllStringSubmatch(remaining, -1) {
		if len(match) > 1 {
			description = strings.Replace(description, match[0], "", -1)
		}
	}
	
	for _, match := range tagRegex.FindAllStringSubmatch(remaining, -1) {
		if len(match) > 0 {
			description = strings.Replace(description, match[0], "", -1)
		}
	}
	
	todo.Description = strings.TrimSpace(description)
	
	return todo, nil
}

func ParseTodos(lines []string) ([]*Todo, error) {
	var todos []*Todo
	for i, line := range lines {
		todo, err := ParseTodo(line)
		if err != nil {
			return nil, err
		}
		if todo != nil {
			todo.ID = i + 1
			todos = append(todos, todo)
		}
	}
	return todos, nil
}