package handler

import (
	"github.com/NeilPanic/SkillsRock_test/internal/repo"
	"github.com/NeilPanic/SkillsRock_test/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct{ svc *service.TaskService }

func Register(r fiber.Router, svc *service.TaskService) {
	h := &TaskHandler{svc}

	r.Post("/", h.create)
	r.Get("/", h.list)
	r.Put("/:id", h.update)
	r.Delete("/:id", h.delete)
}

func (h *TaskHandler) create(c *fiber.Ctx) error {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status,omitempty"` // необязателен
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	if err := h.svc.Create(c.Context(), &task); err != nil {
		return mapErr(err)
	}
	return c.Status(fiber.StatusCreated).JSON(task)
}

func (h *TaskHandler) list(c *fiber.Ctx) error {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	tasks, err := h.svc.List(c.Context(), status, limit, offset)
	if err != nil {
		return mapErr(err)
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var req struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Status      *string `json:"status,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	patch := repo.TaskPatch{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	updated, err := h.svc.Update(c.Context(), id, patch)
	if err != nil {
		return mapErr(err)
	}
	return c.JSON(updated)
}

func (h *TaskHandler) delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return mapErr(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func mapErr(err error) error {
	switch err {
	case service.ErrEmptyTitle, service.ErrInvalidStatus:
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case service.ErrTaskNotFound:
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case service.ErrStatusForbidden:
		return fiber.NewError(fiber.StatusConflict, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
