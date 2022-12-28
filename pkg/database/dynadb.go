package database

import (
	"errors"

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

func NewDynaDB(session *session.Session) (error, *DynaDB) {
	dynaClient := dynamodb.New(session)
	return nil, &DynaDB{session: session, client: dynaClient, tableName: *aws.String("todos")}
}

func (db DynaDB) GetTodos() (error, []model.ToDo) {
	todos := []model.ToDo{}

	input := &dynamodb.ScanInput{TableName: &db.tableName}
	result, err := db.client.Scan(input)
	if err != nil {
		return err, todos
	}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		return err, todos
	}

	return nil, todos
}

func (db DynaDB) SetTodo(todo model.ToDo) error {
	av, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(db.tableName),
	}
	_, err = db.client.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (db DynaDB) GetTodo(id string) (error, model.ToDo) {
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
		return err, todo
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &todo)
	if err != nil {
		return err, todo
	}
	if todo.Id == "" {
		return errors.New("No todo with the specified id found"), todo
	}
	return nil, todo
}

func (db DynaDB) UpdateTodo(id string, todo model.ToDo) error {
	if todo.Id != id { // recreate if id is not identical
		ret := db.DeleteToDo(id)
		if ret != nil {
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
			":title":       {S: &todo.Title},
			":description": {S: &todo.Description},
		},
	}
	_, err := db.client.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (db DynaDB) DeleteToDo(id string) error {
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
		return err
	}
	return nil
}
