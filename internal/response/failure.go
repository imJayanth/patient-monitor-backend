package response

import (
	"net/http"
	"patient-monitor-backend/internal/errors"
	"patient-monitor-backend/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func RespondError(c *fiber.Ctx, err errors.RestAPIError) error {
	return c.Status(err.Status).JSON(NewResponse(nil, err))
}

func RespondUnAuthorised(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusUnauthorized).JSON(NewResponse(nil, errors.NewUnAuthorizedError(helpers.ApplyDefaultArg(message, "Please Login to access the resource"))))
}

func RespondForbidden(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusForbidden).JSON(NewResponse(nil, errors.NewForbiddenError(helpers.ApplyDefaultArg(message, "User not authorized to access this resource"))))
}

func RespondDuplicate(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusConflict).JSON(NewResponse(nil, errors.NewDuplicateRecord(helpers.ApplyDefaultArg(message, "Duplicate record found in database"))))
}

func RespondBadRequest(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusBadRequest).JSON(NewResponse(nil, errors.NewBadRequestError(helpers.ApplyDefaultArg(message, "Request Parameters are not good"))))
}

func RespondNotFound(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusNotFound).JSON(NewResponse([]interface{}{}, errors.NO_ERROR()))
}

func RespondInternalServerError(c *fiber.Ctx, message ...string) error {
	return c.Status(http.StatusInternalServerError).JSON(NewResponse(nil, errors.NewInternalServerError(helpers.ApplyDefaultArg(message, "Something happened unexpectedly"))))
}
