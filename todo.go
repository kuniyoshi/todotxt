package main

import (
	"fmt"
	"strings"
	"time"
)

type Priority rune

const (
	PriorityNone Priority = 0
	PriorityA    Priority = 'A'
	PriorityB    Priority = 'B'
	PriorityC    Priority = 'C'
	PriorityD    Priority = 'D'
	PriorityE    Priority = 'E'
	PriorityF    Priority = 'F'
	PriorityG    Priority = 'G'
	PriorityH    Priority = 'H'
	PriorityI    Priority = 'I'
	PriorityJ    Priority = 'J'
	PriorityK    Priority = 'K'
	PriorityL    Priority = 'L'
	PriorityM    Priority = 'M'
	PriorityN    Priority = 'N'
	PriorityO    Priority = 'O'
	PriorityP    Priority = 'P'
	PriorityQ    Priority = 'Q'
	PriorityR    Priority = 'R'
	PriorityS    Priority = 'S'
	PriorityT    Priority = 'T'
	PriorityU    Priority = 'U'
	PriorityV    Priority = 'V'
	PriorityW    Priority = 'W'
	PriorityX    Priority = 'X'
	PriorityY    Priority = 'Y'
	PriorityZ    Priority = 'Z'
)

type Todo struct {
	ID             int
	Complete       bool
	Priority       Priority
	CreationDate   *time.Time
	CompletionDate *time.Time
	Description    string
	Projects       []string
	Contexts       []string
	Tags           map[string]string
	Raw            string
}

func NewTodo(description string) *Todo {
	now := time.Now()
	return &Todo{
		Complete:     false,
		Priority:     PriorityNone,
		CreationDate: &now,
		Description:  description,
		Projects:     []string{},
		Contexts:     []string{},
		Tags:         make(map[string]string),
		Raw:          description,
	}
}

func (t *Todo) String() string {
	var parts []string

	if t.Complete {
		parts = append(parts, "x")
		if t.CompletionDate != nil {
			parts = append(parts, t.CompletionDate.Format("2006-01-02"))
		}
		if t.CreationDate != nil {
			parts = append(parts, t.CreationDate.Format("2006-01-02"))
		}
	} else {
		if t.Priority != PriorityNone {
			parts = append(parts, fmt.Sprintf("(%c)", t.Priority))
		}
		if t.CreationDate != nil {
			parts = append(parts, t.CreationDate.Format("2006-01-02"))
		}
	}

	// Build full description with projects, contexts, and tags
	fullDesc := t.Description

	for _, project := range t.Projects {
		if !strings.Contains(fullDesc, "+"+project) {
			fullDesc += " +" + project
		}
	}

	for _, context := range t.Contexts {
		if !strings.Contains(fullDesc, "@"+context) {
			fullDesc += " @" + context
		}
	}

	for key, value := range t.Tags {
		tag := fmt.Sprintf("%s:%s", key, value)
		if !strings.Contains(fullDesc, tag) {
			fullDesc += " " + tag
		}
	}

	parts = append(parts, fullDesc)

	return strings.Join(parts, " ")
}

func (t *Todo) SetPriority(p Priority) {
	if !t.Complete {
		t.Priority = p
	}
}

func (t *Todo) MarkComplete() {
	t.Complete = true
	now := time.Now()
	t.CompletionDate = &now
	t.Priority = PriorityNone
}

func (t *Todo) MarkUncomplete() {
	t.Complete = false
	t.CompletionDate = nil
}

func (t *Todo) AddProject(project string) {
	for _, p := range t.Projects {
		if p == project {
			return
		}
	}
	t.Projects = append(t.Projects, project)
}

func (t *Todo) AddContext(context string) {
	for _, c := range t.Contexts {
		if c == context {
			return
		}
	}
	t.Contexts = append(t.Contexts, context)
}

func (t *Todo) AddTag(key, value string) {
	t.Tags[key] = value
}

func (t *Todo) GetDueDate() *time.Time {
	if due, ok := t.Tags["due"]; ok {
		if date, err := time.Parse("2006-01-02", due); err == nil {
			return &date
		}
	}
	return nil
}
