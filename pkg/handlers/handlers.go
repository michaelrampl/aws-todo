package handlers

import (
	"github.com/michaelrampl/aws-todo/pkg/database"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

func V1TodoGet(db database.Database) (error, []model.ToDo) {
	return db.GetTodos()
}

func V1TodoPut(db database.Database, todo model.ToDo) error {
	return db.SetTodo(todo)
}

func V1TodoGetByID(db database.Database, id string) (error, model.ToDo) {
	return db.GetTodo(id)
}

func V1TodoPutByID(db database.Database, id string, todo model.ToDo) error {
	return db.UpdateTodo(id, todo)
}

func V1TodoDeleteByID(db database.Database, id string) error {
	return db.DeleteToDo(id)
}
