package main

import (
	"encoding/json"
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
	db database.DynaDB
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	db = database.NewDynaDB(awsSession)
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

func statusResponse(status int) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status
	return &resp, nil
}

func handler(req events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	id, hasID := req.PathParameters["id"]

	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = 200

	switch req.RequestContext.HTTP.Method {
	case "GET":
		if hasID {
			var status, data = handlers.V1TodoGetByID(db, id)
			if status == 200 {
				return jsonResponse(status, data)
			} else {
				return statusResponse(status)
			}
		} else {
			var status, data = handlers.V1TodoGet(db)
			if status == 200 {
				return jsonResponse(status, data)
			} else {
				return statusResponse(status)
			}
		}
	case "PUT":
		if hasID {
			todo := model.ToDo{}
			if err := json.Unmarshal([]byte(req.Body), &todo); err != nil {
				return statusResponse(400)
			}
			return statusResponse(handlers.V1TodoPutByID(db, id, todo))
		} else {
			todo := model.ToDo{}
			if err := json.Unmarshal([]byte(req.Body), &todo); err != nil {
				return statusResponse(400)
			}
			return statusResponse(handlers.V1TodoPut(db, todo))
		}
	case "DELETE":
		if hasID {
			return statusResponse(handlers.V1TodoDeleteByID(db, id))
		} else {
			return statusResponse(400)
		}
	default:
		return statusResponse(499)
	}
}
