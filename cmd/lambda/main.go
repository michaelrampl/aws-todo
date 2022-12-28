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

// Global variable storing the dynamo db instance
var (
	db *database.DynaDB
)

func main() {
	// Create the aws session object, using the region info stored as environment variable
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		// Stop Application if there has been an error creating the aws session
		log.Fatalf("Could not create aws session: %s", err)
	}

	// Connect to the Dynamo DB
	err, db = database.NewDynaDB(awsSession)

	// Start the handler
	lambda.Start(handler)
}

// Generates a json response based on an http status and a serializable object
func jsonResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {

	// set the json header
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}

	// set the http status
	resp.StatusCode = status

	// serialize the json object and set the body
	jsonBody, _ := json.Marshal(body)
	resp.Body = string(jsonBody)

	// return the response object
	return &resp, nil
}

// Creates an http success response containing data
func dataResponse(data interface{}) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusOK, data)
}

// Creates an http success response containing a success message
func messageResponse(msg string) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusOK, model.NewSuccessMessage(msg))
}

// Creates an http success response containing an error message
func errorResponse(msg string) (*events.APIGatewayProxyResponse, error) {
	return jsonResponse(http.StatusBadRequest, model.NewErrorMessage(msg))
}

// Main handler, being run for every route
func handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	// detect if the id parameter has been sent
	id, hasID := req.PathParameters["id"]

	// Decide what todo based on the HTTP Method (Get, Put, Delete)
	switch req.RequestContext.HTTP.Method {
	case "GET":
		if hasID { // If an id has been provided, get a single todo object based on this id
			err, data := handlers.V1TodoGetByID(db, id)
			if err != nil {
				log.Printf("Error while getting todo %s: %s", id, err)
				return errorResponse("There was an error loading the To-Do object.")
			} else {
				return dataResponse(data)
			}
		} else { // If no id has been provided, get all todo objects
			err, data := handlers.V1TodoGet(db)
			if err != nil {
				log.Printf("Error while getting todos: %s", err)
				return errorResponse("There was an error loading the To-Do objects.")
			} else {
				return dataResponse(data)
			}
		}
	case "PUT":
		if hasID { // If an id has been provided, update an existing todo object based on the id
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
		} else { // if no id has been provided, create a new todo object
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
		if hasID { // If an id has been provided, delete a signle todo object based on the id
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
