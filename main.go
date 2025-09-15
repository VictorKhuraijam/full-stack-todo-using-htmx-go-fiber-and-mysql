package main

import (
	"fmt"
	"log"
	"todo/server/config"
	"todo/server/database"
	"todo/server/handlers"
	"todo/server/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func main() {
	//Load configuration
	cfg := config.Load() //Load DB/Server configs

	//Connect to database
	db, err := database.Connect(cfg) //Connects to MySQL
	if err != nil{
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// check db conectivity
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to MySQL!")

	//Initialize template engine
	engine := html.New("./client/templates", ".html")//Initializes HTML rendering engine
	engine.Reload(true) //enable live reloads(useful during development)

	//Create Fiber app
	app := fiber.New(fiber.Config{ // 	Creates HTTP server
		Views: engine, //to render HTML templates
	})



	//Middleware
	//Adds a request logger middleware that prints method, path, status, and response time to the terminal for every HTTP request.
	app.Use(logger.New())

	//Static files
	app.Static("/static", "./client/static")

	//Initialize services and handlers
	todoService := services.NewTodoService(db) //handles database operations (logic layer).
	todoHandler := handlers.NewTodoHandler(todoService) // handles incoming HTTP requests and uses todoService to perform actual operations.

	//Routes
	app.Get("/", todoHandler.Index)
	// app.Get("/todos/todo/:id", todoHandler.GetTodo)
	app.Post("/todos", todoHandler.Create)
	app.Get("/todos/:id", todoHandler.EditForm)
	app.Put("/todos/toggle/:id", todoHandler.Toggle)
	app.Put("/todos/update/:id", todoHandler.Update)
	app.Delete("/todos/:id", todoHandler.Delete)

	//Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port)) // fiber http server starts
}
