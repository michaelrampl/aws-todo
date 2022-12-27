package database

import (
	"github.com/michaelrampl/aws-todo/pkg/model"
)

type Database interface {
	GetTodos() (int, []model.ToDo)
	SetTodo(todo model.ToDo) int
	UpdateTodo(id string, todo model.ToDo) int
	GetTodo(id string) (int, model.ToDo)
	DeleteToDo(id string) int
}
