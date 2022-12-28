package handlers

import (
	"github.com/michaelrampl/aws-todo/pkg/database"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

// Gets all todo objects stored in the databse
func V1TodoGet(db database.Database) (error, []model.ToDo) {
	return db.GetTodos()
}

// Creates a single todo object
func V1TodoPut(db database.Database, todo model.ToDo) error {
	return db.SetTodo(todo)
}

// Gets a single todo object by id
func V1TodoGetByID(db database.Database, id string) (error, model.ToDo) {
	return db.GetTodo(id)
}

// Updates a single todo object based on its id
func V1TodoPutByID(db database.Database, id string, todo model.ToDo) error {
	return db.UpdateTodo(id, todo)
}

// Deletes a single todo object based on its id
func V1TodoDeleteByID(db database.Database, id string) error {
	return db.DeleteToDo(id)
}
