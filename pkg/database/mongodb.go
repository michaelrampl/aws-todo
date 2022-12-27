package database

import (
	"context"
	"log"
	"time"

	"github.com/michaelrampl/aws-todo/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client    *mongo.Client
	context   *context.Context
	tableName string
}

func NewMongoDB(uri string) MongoDB {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return MongoDB{client: client, context: &ctx, tableName: "todos"}
}

func (db MongoDB) GetTodos() (int, []model.ToDo) {
	log.Printf("GetTodos")
	todos := []model.ToDo{}
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	cursor, err := collection.Find(*db.context, bson.M{})
	if err != nil {
		return 400, todos
	}
	if err = cursor.All(*db.context, &todos); err != nil {
		log.Printf("Error while loading all todos: %s", err)
		return 400, todos
	}
	return 200, todos

}

func (db MongoDB) SetTodo(todo model.ToDo) int {
	log.Printf("SetTodo")
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.InsertOne(*db.context, todo)
	if err != nil {
		log.Printf("Error while inserting todo into the database: %s", err)
		return 400
	}
	return 200
}

func (db MongoDB) GetTodo(id string) (int, model.ToDo) {
	log.Printf("GetTodo")
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	todo := model.ToDo{}
	err := collection.FindOne(*db.context, bson.M{"id": id}).Decode(&todo)
	if err != nil {
		log.Printf("Error while loading todo with id %s from the database: %s", id, err)
		return 400, todo
	}
	return 200, todo
}

func (db MongoDB) UpdateTodo(id string, todo model.ToDo) int {
	log.Printf("UpdateTodo")
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.UpdateOne(*db.context, bson.M{"id": id}, bson.D{{"$set", todo}})
	if err != nil {
		log.Printf("Error while updating todo with id %s from the database: %s", id, err)
		return 400
	}
	return 200
}

func (db MongoDB) DeleteToDo(id string) int {
	log.Printf("DeleteToDo")
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.DeleteOne(*db.context, bson.M{"id": id})
	if err != nil {
		log.Printf("Error while deleting todo with id %s from the database: %s", id, err)
		return 400
	}
	return 200
}
