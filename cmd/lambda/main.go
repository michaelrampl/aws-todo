package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/michaelrampl/aws-todo/pkg/database"
	"github.com/michaelrampl/aws-todo/pkg/handlers"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

var (
	db *database.DynaDB
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	err, db = database.NewDynaDB(awsSession)
	lambda.Start(handler)
}

const tableName = "LambdaInGoUser"

func jsonResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status
	jsonBody, _ := json.Marshal(body)
	resp.Body = string(jsonBody)
	return &resp, nil
}

func dataResponse(data interface{}) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusOK, data)
}

func messageResponse(msg string) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusOK, model.NewSuccessMessage(msg))
}

func errorResponse(msg string) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusBadRequest, model.NewErrorMessage(msg))
}

func handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	id, hasID := req.PathParameters["id"]

	switch req.RequestContext.HTTP.Method {
	case "GET":
		if hasID {
			err, data := handlers.V1TodoGetByID(db, id)
			if err != nil {
				log.Printf("Error while getting todo %s: %s", id, err)
				return errorResponse("There was an error loading the To-Do object.")
			} else {
				return dataResponse(data)
			}
		} else {
			err, data := handlers.V1TodoGet(db)
			if err != nil {
				log.Printf("Error while getting todos: %s", err)
				return errorResponse("There was an error loading the To-Do objects.")
			} else {
				return dataResponse(data)
			}
		}
	case "PUT":
		if hasID {
			todo := model.ToDo{}
			if err := json.Unmarshal([]byte(req.Body), &todo); err != nil {
				log.Printf("Error putting existing todo %s: %s", id, err)
				return errorResponse("There was an error while updating the To-Do object.")
			}
			err := handlers.V1TodoPutByID(db, id, todo)
			if err != nil {
				log.Printf("Error while updating todo %s: %s", id, err)
				return errorResponse("There was an error while updating the To-Do object.")
			} else {
				return messageResponse("To-Do object updated successfully.")
			}
		} else {
			todo := model.ToDo{}
			if err := json.Unmarshal([]byte(req.Body), &todo); err != nil {
				log.Printf("Error while creating todo: %s", err)
				return errorResponse("There was an error while creating a new To-Do object.")
			}
			err := handlers.V1TodoPut(db, todo)
			if err != nil {
				log.Printf("Error while creating todo: %s", err)
				return errorResponse("There was an error while creating a new To-Do object.")
			} else {
				return messageResponse("To-Do object created successfully.")
			}
		}
	case "DELETE":
		if hasID {
			err := handlers.V1TodoDeleteByID(db, id)
			if err != nil {
				log.Printf("Error while deleting todo %s: %s", id, err)
				return errorResponse("There has been an error while deleting the To-Do object.")
			} else {
				return messageResponse("To-Do object deleted successfully.")
			}
		} else {
			return jsonResponse(http.StatusNotFound, model.NewErrorMessage("Invalid Route."))
		}
	default:
		return jsonResponse(http.StatusNotFound, model.NewErrorMessage("Invalid Route."))
	}
}
