package handler

import (
	"github.com/NeilPanic/SkillsRock_test/todo-app/internal/repo"
	"strconv"

	"SkillsRock_test/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct{ svc *service.TaskService }

func Register(app fiber.Router, svc *service.TaskService) {
	h := &TaskHandler{svc}
	app.Post("/", h.create)
	app.Get("/", h.list)
	// PUT, DELETE…
}

func (h *TaskHandler) create(c *fiber.Ctx) error {
	var req struct {
		Title, Description string `json:"title"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	task := repo.Task{Title: req.Title, Description: req.Description}
	if err := h.svc.Create(c.Context(), &task); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskHandler) list(c *fiber.Ctx) error {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	tasks, err := h.svc.List(c.Context(), status, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(tasks)
}
