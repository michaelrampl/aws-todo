package database

import (
	"github.com/michaelrampl/aws-todo/pkg/model"
)

// Interface describing the potential database access methods
// The classes implementing this interface will then be used by the main
// binaries handling the connections and injected into the handlers
type Database interface {

	// Get all todos
	GetTodos() (error, []model.ToDo)

	// Set a single todo object
	SetTodo(todo model.ToDo) error

	// Update a single todo object based on its id
	UpdateTodo(id string, todo model.ToDo) error

	// Get a single todo object based on its id
	GetTodo(id string) (error, model.ToDo)

	// Delete a single todo object based on its id
	DeleteToDo(id string) error
}
