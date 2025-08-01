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
	cfg := config.Load()

	//Connect to database
	db, err := database.Connect(cfg)
	if err != nil{
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to MySQL!")

	//Initialize template engine
	engine := html.New("./client/templates", ".html")
	engine.Reload(true)

	//Create Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//Middleware
	app.Use(logger.New())

	//Static files
	app.Static("/static", "./client/static")

	//Initialize services and handlers
	todoService := services.NewTodoService(db)
	todoHandler := handlers.NewTodoHandler(todoService)

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
	log.Fatal(app.Listen(":" + cfg.Port))
}
