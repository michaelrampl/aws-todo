package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

type DynaDB struct {
	session   *session.Session
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynaDB(session *session.Session) DynaDB {
	dynaClient := dynamodb.New(session)
	return DynaDB{session: session, client: dynaClient, tableName: *aws.String("todos")}
}

func (db DynaDB) GetTodos() (int, []model.ToDo) {
	todos := []model.ToDo{}

	input := &dynamodb.ScanInput{TableName: &db.tableName}
	result, err := db.client.Scan(input)
	if err != nil {
		log.Printf("Error while loading all todos: %s", err)
		return 400, todos

	}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		log.Printf("Error while unmarshalling todos: %s", err)
		return 400, todos
	}

	return 200, todos

}

func (db DynaDB) SetTodo(todo model.ToDo) int {
	av, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		log.Printf("Error while marshalling todo: %s", err)
		return 400
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.tableName),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		log.Printf("Error while saving todo: %s", err)
		return 400
	}

	return 200
}

func (db DynaDB) GetTodo(id string) (int, model.ToDo) {
	todo := model.ToDo{}

	input := &dynamodb.GetItemInput{
		TableName: &db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	result, err := db.client.GetItem(input)
	if err != nil {
		log.Printf("Error while loading todo %s: %s", id, err)
		return 400, todo
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &todo)
	if err != nil {
		log.Printf("Error while unmarshalling todo %s: %s", id, err)
		return 400, todo
	}
	if todo.Id == "" {
		log.Printf("Could not find todo %s", id)
		return 400, todo
	}

	return 200, todo
}

func (db DynaDB) UpdateTodo(id string, todo model.ToDo) int {
	if todo.Id != id { // recreate if id is not identical
		ret := db.DeleteToDo(id)
		if ret != 200 {
			return ret
		}
		return db.SetTodo(todo)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: &db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression: aws.String("set title = :title, description = :description"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			//":id":          {S: &todo.Id},
			":title":       {S: &todo.Title},
			":description": {S: &todo.Description},
		},
	}

	_, err := db.client.UpdateItem(input)
	if err != nil {
		log.Printf("Error while updating todo %s: %s", id, err)
		return 400
	}

	return 200
}

func (db DynaDB) DeleteToDo(id string) int {
	input := &dynamodb.DeleteItemInput{
		TableName: &db.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}
	_, err := db.client.DeleteItem(input)
	if err != nil {
		log.Printf("Error while deleting todo %s: %s", id, err)
		return 400
	}

	return 200
}
