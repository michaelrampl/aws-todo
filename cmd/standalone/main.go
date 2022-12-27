package main

import (
	"os"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelrampl/aws-todo/pkg/database"
	"github.com/michaelrampl/aws-todo/pkg/handlers"
	"github.com/michaelrampl/aws-todo/pkg/model"
)

func main() {
	app := fiber.New()
	db_uri := os.Getenv("TODO_MONGODB_URI")
	db := database.NewMongoDB(db_uri)

	app.Get("/v1/todo", func(c *fiber.Ctx) error {
		var status, data = handlers.V1TodoGet(db)
		if status == 200 {
			return c.Status(status).JSON(data)
		} else {
			return c.SendStatus(status)
		}
	})

	app.Put("/v1/todo", func(c *fiber.Ctx) error {
		todo := model.ToDo{}
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error on Route [%s|%s] while parsing body: %s", c.Route().Method, c.Route().Path, err)
			return c.SendStatus(400)
		}
		return c.SendStatus(handlers.V1TodoPut(db, todo))
	})

	app.Get("/v1/todo/:id", func(c *fiber.Ctx) error {
		var status, data = handlers.V1TodoGetByID(db, c.Params("id"))
		if status == 200 {
			return c.Status(status).JSON(data)
		} else {
			return c.SendStatus(status)
		}
	})

	app.Put("/v1/todo/:id", func(c *fiber.Ctx) error {
		todo := model.ToDo{}
		if err := c.BodyParser(&todo); err != nil {
			log.Printf("Error on Route [%s|%s] while parsing body: %s", c.Route().Method, c.Route().Path, err)
			return c.SendStatus(400)
		}
		return c.SendStatus(handlers.V1TodoPutByID(db, c.Params("id"), todo))
	})

	app.Delete("/v1/todo/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(handlers.V1TodoDeleteByID(db, c.Params("id")))
	})

	app.Listen(":3000")
}
