package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelrampl/aws-todo/pkg/database"
	"github.com/michaelrampl/aws-todo/pkg/globals"
	"github.com/michaelrampl/aws-todo/pkg/handlers"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

func main() {

	// Create instance of go fiber app
	app := fiber.New()

	// Create the instance of a MongoDB connector
	// Load the connection uri from the environment
	err, db := database.NewMongoDB(os.Getenv("TODO_MONGODB_URI"))
	if err != nil {
		// Stop Application if there has been an error connecting to the db
		log.Fatalf("Could not connect to mongodb: %s", err)
	}

	// Get all todo objects from the database
	app.Get("/v1/todo", func(c *fiber.Ctx) error {
		err, data := handlers.V1TodoGet(db)
		if err != nil {
			log.Printf("Error while getting todos: %s", err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_GET_ALL_ERROR))
		} else {
			return c.Status(http.StatusOK).JSON(data)
		}
	})

	// Put a new todo object
	app.Put("/v1/todo", func(c *fiber.Ctx) error {
		todo := model.ToDo{}
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error while parsing body in route %s: %s", c.Route().Path, err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_CREATE_ERROR))
		}

		err := handlers.V1TodoPut(db, todo)
		if err != nil {
			log.Printf("Error while creating todo %s: %s", c.Params("id"), err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_CREATE_ERROR))
		} else {
			return c.Status(http.StatusOK).JSON(model.NewSuccessMessage(globals.TODO_CREATE_SUCCESS))
		}

	})

	// Get a single todo object based on the id
	app.Get("/v1/todo/:id", func(c *fiber.Ctx) error {
		err, data := handlers.V1TodoGetByID(db, c.Params("id"))
		if err != nil {
			log.Printf("Error while getting todo %s: %s", c.Params("id"), err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_GET_ERROR))
		} else {
			return c.Status(http.StatusOK).JSON(data)
		}
	})

	// Put/Update an existing todo object based on its id
	app.Put("/v1/todo/:id", func(c *fiber.Ctx) error {
		todo := model.ToDo{}
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error while parsing body in route %s: %s", c.Route().Path, err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_UPDATE_ERROR))
		}
		err := handlers.V1TodoPutByID(db, c.Params("id"), todo)
		if err != nil {
			log.Printf("Error while updating todo %s: %s", c.Params("id"), err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_UPDATE_ERROR))
		} else {
			return c.Status(http.StatusOK).JSON(model.NewSuccessMessage(globals.TODO_UPDATE_SUCCESS))
		}
	})

	// Delete an existing todo object based on its id
	app.Delete("/v1/todo/:id", func(c *fiber.Ctx) error {
		err := handlers.V1TodoDeleteByID(db, c.Params("id"))
		if err != nil {
			log.Printf("Error deleting getting todo %s: %s", c.Params("id"), err)
			return c.Status(http.StatusBadRequest).JSON(model.NewErrorMessage(globals.TODO_DELETE_ERROR))
		} else {
			return c.Status(http.StatusOK).JSON(model.NewSuccessMessage(globals.TODO_DELETE_SUCCESS))
		}
	})

	app.Listen(":3000")
}
