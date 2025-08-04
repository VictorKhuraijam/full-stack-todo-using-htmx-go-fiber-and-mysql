package handlers

import (
	"log" // Used for logging errors or debugging info.
	"strconv" //Used for converting string to int (like id params).
	"todo/server/services"

	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) Index(c *fiber.Ctx) error  {
	todos,err := h.service.GetAll()
	if err != nil{
		log.Println(err)
		return c.Status(500).SendString("Error fetching todos")
	}

	return c.Render("index", fiber.Map{
		"Todos": todos,
		"Empty": len(todos) == 0,
	})
}

func (h *TodoHandler) Create(c *fiber.Ctx) error  {
	title := c.FormValue("title")
	if title == ""{
		return c.Status(400).SendString("Title is required")
	}

	todo, err := h.service.Create(title)
	if err != nil{
		return c.Status(500).SendString("Error creating todo")
	}

	return c.Render("partials/todo-item",todo)
}

func (h *TodoHandler) GetTodo(c *fiber.Ctx) error  {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid Id")
	}

	todo, err := h.service.GetByID(id)
	if err != nil{
		return c.Status(500).SendString("Error getting todo")
	}

	log.Printf("get todo: %+v\n", todo)

	return c.Render("partials/todo-edit", todo)
}

func (h *TodoHandler) Toggle(c *fiber.Ctx) error  {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.Status(400).SendString("Invalid ID")
	}

	todo, err := h.service.ToggleComplete(id)
	if err != nil{
		return c.Status(500).SendString("Error updating toggle todo")
	}

	log.Printf("Toggled todo: %+v\n", todo)

	return c.Render("partials/todo-item", todo)
}

func (h *TodoHandler) EditForm(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	todo, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(500).SendString("Todo not found")
	}

	mode := c.Query("mode")

	if mode == "view" {
		// render the read-only todo-item view
		return c.Render("partials/todo-item", todo)
	}

	return c.Render("partials/todo-edit", todo)
}


func (h *TodoHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("Invalid ID:", err)
		return c.Status(400).SendString("Invalid ID")
	}

	title := c.FormValue("title")
	log.Println("Title:", title)
	if title == ""{
		return c.Status(400).SendString("Title is required")
	}

	completedStr := c.FormValue("completed")
	log.Println("Completed (raw):", completedStr)

	// completed, err := strconv.ParseBool(completedStr)
	// if err != nil{
	// 	return c.Status(400).SendString("Invalid value for completed")
	// }
	completed := completedStr == "true"
	log.Println("Completed (parsed):", completed)


	todo, err := h.service.UpdateTodo(id, title, completed)
	if err != nil {
		log.Println("Update failed:", err)
		return c.Status(500).SendString("Error updating Todo")
	}

	log.Printf("Updated todo: %+v\n", todo)
	return c.Render("partials/todo-item", todo)
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error  {
	id,err := strconv.Atoi(c.Params("id")) //Converts the URL parameter id from string to int.
	if err != nil{
		return c.Status(400).SendString("Invalid ID")
	}

	err = h.service.Delete(id)
	if err != nil{
		return c.Status(500).SendString("Error deleting todo")
	}

	// return c.SendString("Todo Deleted")
	return nil
}
