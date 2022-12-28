package model

import (
	"gopkg.in/go-playground/validator.v9"
)

// Represents a todo object
type ToDo struct {

	// Identifier of the ToDo object in uuid format
	Id string `json:"id" bson:"id" validate:"required,uuid4"`

	// Title of the ToDo object
	Title string `json:"title" bson:"title" validate:"required,ascii,min=3,max=100"`

	// Optional description of the ToDo object
	Description string `json:"description" bson:"description" validate:"ascii,max=500"`
}

// Checks if an ToDO object is valid
func (todo *ToDo) Validate() error {
	validate := validator.New()
	return validate.Struct(todo)
}

// Represents an Success Message returned to the user
type SuccessMessage struct {

	// The message shown to the user
	Message string `json:"message"`
}

// Represents an Error Message returned to the user
type ErrorMessage struct {

	// The message shown to the user
	Error string `json:"error"`
}

// Constructs a new ErrorMessage object
func NewErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{
		Error: msg,
	}
}

// Constructs a new SuccessMessage object
func NewSuccessMessage(msg string) SuccessMessage {
	return SuccessMessage{
		Message: msg,
	}
}
