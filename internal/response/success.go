package response

import (
	"net/http"
	"patient-monitor-backend/internal/errors"

	"github.com/gofiber/fiber/v2"
)

func RespondSuccess(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(NewResponse(nil, errors.NO_ERROR()))
}

func RespondStringSuccess(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusOK).JSON(NewResponse(message, errors.NO_ERROR()))
}

func RespondStringArraySuccess(c *fiber.Ctx, messages []string) error {
	return c.Status(http.StatusOK).JSON(NewResponse(messages, errors.NO_ERROR()))
}

func RespondCreated(c *fiber.Ctx, obj interface{}) error {
	return c.Status(http.StatusCreated).JSON(NewResponse(obj, errors.NO_ERROR()))
}

func RespondSuccessJSON(c *fiber.Ctx, obj interface{}) error {
	return c.Status(http.StatusOK).JSON(NewResponse(obj, errors.NO_ERROR()))
}
