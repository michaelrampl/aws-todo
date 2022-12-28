package database

import (
	"context"
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

func NewMongoDB(uri string) (error, *MongoDB) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err, nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	err = client.Connect(ctx)
	if err != nil {
		return err, nil
	}
	return nil, &MongoDB{client: client, context: &ctx, tableName: "todos"}
}

func (db MongoDB) GetTodos() (error, []model.ToDo) {
	todos := []model.ToDo{}
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	cursor, err := collection.Find(*db.context, bson.M{})
	if err != nil {
		return err, todos
	}
	if err = cursor.All(*db.context, &todos); err != nil {
		return err, todos
	}
	return nil, todos

}

func (db MongoDB) SetTodo(todo model.ToDo) error {
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.InsertOne(*db.context, todo)
	if err != nil {
		return err
	}
	return nil
}

func (db MongoDB) GetTodo(id string) (error, model.ToDo) {
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	todo := model.ToDo{}
	err := collection.FindOne(*db.context, bson.M{"id": id}).Decode(&todo)
	if err != nil {
		return err, todo
	}
	return nil, todo
}

func (db MongoDB) UpdateTodo(id string, todo model.ToDo) error {
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.UpdateOne(*db.context, bson.M{"id": id}, bson.D{{"$set", todo}})
	if err != nil {
		return err
	}
	return nil
}

func (db MongoDB) DeleteToDo(id string) error {
	collection := db.client.Database(db.tableName).Collection(db.tableName)
	_, err := collection.DeleteOne(*db.context, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
