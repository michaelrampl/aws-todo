package model

type ToDo struct {
	Id          string `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func NewErrorMessage(msg string) ErrorMessage {
	return ErrorMessage{
		Error: msg,
	}
}

func NewSuccessMessage(msg string) SuccessMessage {
	return SuccessMessage{
		Message: msg,
	}
}
