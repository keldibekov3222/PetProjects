package main

import (
	"fiber-go/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()
	api := app.Group("/employee")
	api.Get("/", handlers.GetAllEmployees)
	api.Get("/:id", handlers.GetEmployeeByID)
	api.Post("/", handlers.CreateNewEmployee)
	api.Put("/:id", handlers.UpdateEmployee)
	api.Delete("/:id", handlers.DeleteEmployee)

	log.Fatal(app.Listen(":8080"))
}
