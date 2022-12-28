package database

import (
	"github.com/michaelrampl/aws-todo/pkg/model"
)

type Database interface {
	GetTodos() (error, []model.ToDo)
	SetTodo(todo model.ToDo) error
	UpdateTodo(id string, todo model.ToDo) error
	GetTodo(id string) (error, model.ToDo)
	DeleteToDo(id string) error
}
