package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TodoFile struct {
	Path  string
	Todos []*Todo
}

func NewTodoFile(path string) *TodoFile {
	return &TodoFile{
		Path:  path,
		Todos: []*Todo{},
	}
}

func (tf *TodoFile) Load() error {
	if _, err := os.Stat(tf.Path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(tf.Path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	todos, err := ParseTodos(lines)
	if err != nil {
		return fmt.Errorf("failed to parse todos: %w", err)
	}

	tf.Todos = todos
	return nil
}

func (tf *TodoFile) Save() error {
	dir := filepath.Dir(tf.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(tf.Path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, todo := range tf.Todos {
		if _, err := writer.WriteString(todo.String() + "\n"); err != nil {
			return fmt.Errorf("failed to write todo: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func (tf *TodoFile) Add(todo *Todo) {
	todo.ID = len(tf.Todos) + 1
	tf.Todos = append(tf.Todos, todo)
}

func (tf *TodoFile) GetByID(id int) *Todo {
	for _, todo := range tf.Todos {
		if todo.ID == id {
			return todo
		}
	}
	return nil
}

func (tf *TodoFile) Delete(id int) bool {
	for i, todo := range tf.Todos {
		if todo.ID == id {
			tf.Todos = append(tf.Todos[:i], tf.Todos[i+1:]...)
			tf.reindexTodos()
			return true
		}
	}
	return false
}

func (tf *TodoFile) reindexTodos() {
	for i, todo := range tf.Todos {
		todo.ID = i + 1
	}
}

func (tf *TodoFile) Search(query string) []*Todo {
	query = strings.ToLower(query)
	var results []*Todo

	for _, todo := range tf.Todos {
		if strings.Contains(strings.ToLower(todo.Description), query) {
			results = append(results, todo)
			continue
		}

		for _, project := range todo.Projects {
			if strings.Contains(strings.ToLower(project), query) {
				results = append(results, todo)
				break
			}
		}

		for _, context := range todo.Contexts {
			if strings.Contains(strings.ToLower(context), query) {
				results = append(results, todo)
				break
			}
		}
	}

	return results
}

func (tf *TodoFile) FilterByProject(project string) []*Todo {
	var results []*Todo
	for _, todo := range tf.Todos {
		for _, p := range todo.Projects {
			if p == project {
				results = append(results, todo)
				break
			}
		}
	}
	return results
}

func (tf *TodoFile) FilterByContext(context string) []*Todo {
	var results []*Todo
	for _, todo := range tf.Todos {
		for _, c := range todo.Contexts {
			if c == context {
				results = append(results, todo)
				break
			}
		}
	}
	return results
}

func (tf *TodoFile) GetCompleted() []*Todo {
	var results []*Todo
	for _, todo := range tf.Todos {
		if todo.Complete {
			results = append(results, todo)
		}
	}
	return results
}

func (tf *TodoFile) GetIncomplete() []*Todo {
	var results []*Todo
	for _, todo := range tf.Todos {
		if !todo.Complete {
			results = append(results, todo)
		}
	}
	return results
}
