package handlers

import (
	"log"
	"strconv"
	"todo/internal/services"

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

func (h *TodoHandler) Delete(c *fiber.Ctx) error  {
	id,err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.Status(400).SendString("Invalid ID")
	}

	err = h.service.Delete(id)
	if err != nil{
		return c.Status(500).SendString("Error deleting todo")
	}

	return c.SendString("Todo Deleted")
}
