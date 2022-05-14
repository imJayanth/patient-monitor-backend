package controllers

import (
	"patient-monitor-backend/internal/config"
	"patient-monitor-backend/internal/errors"
	"patient-monitor-backend/internal/helpers"
	"patient-monitor-backend/internal/models"
	"patient-monitor-backend/internal/response"
	"patient-monitor-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type UserController struct {
	UserService *services.UserService
	appConf     *config.AppConfigModel
	name        string
	logger      *zap.Logger
}

func NewUserController() *UserController {
	return &UserController{UserService: services.NewUserService(), appConf: config.AppConfig, name: "UserController", logger: config.Logger}
}

func (ac *UserController) Login(c *fiber.Ctx) error {
	m := "Login"
	ac.logInfo(m, "Start", nil)
	defer ac.logInfo(m, "Complete", nil)

	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	ac.logInfo(m, "Input", []zapcore.Field{zap.Any("User", user.ToJsonRaw())})
	if saveErr := ac.UserService.Login(&user); saveErr.IsNotNull() {
		return response.RespondError(c, saveErr)
	}
	return response.RespondStringSuccess(c, "success")
}

func (ac *UserController) AddContact(c *fiber.Ctx) error {
	m := "AddContact"
	ac.logInfo(m, "Start", nil)
	defer ac.logInfo(m, "Complete", nil)

	contact := models.Contact{}
	if err := c.BodyParser(&contact); err != nil {
		return err
	}
	ac.logInfo(m, "Input", []zapcore.Field{zap.Any("contact", contact.ToJsonRaw())})
	if saveErr := ac.UserService.AddContact(&contact); saveErr.IsNotNull() {
		return response.RespondError(c, saveErr)
	}
	return response.RespondSuccessJSON(c, &contact)
}

func (ac *UserController) GetContactsById(c *fiber.Ctx) error {
	m := "GetContantsById"
	ac.logInfo(m, "Start", nil)
	defer ac.logInfo(m, "Complete", nil)

	patientId, e := c.ParamsInt("ID")
	if e != nil {
		err := errors.NewUnAuthorizedError("Invalid patient Id")
		ac.logError(m, "Params", []zapcore.Field{zap.Any("ID", err.Message)})
		return response.RespondError(c, err)
	}

	ac.logInfo(m, "Input", []zapcore.Field{zap.Any("Id", patientId)})
	contacts := []models.Contact{}
	if getErr := ac.UserService.GetContactsById(&contacts, patientId); getErr.IsNotNull() {
		return response.RespondError(c, getErr)
	}
	return response.RespondSuccessJSON(c, &contacts)
}

func (ac *UserController) GetSensorDataById(c *fiber.Ctx) error {
	m := "GetSensorDataById"
	ac.logInfo(m, "Start", nil)
	defer ac.logInfo(m, "Complete", nil)

	patientId, e := c.ParamsInt("ID")
	if e != nil {
		err := errors.NewUnAuthorizedError("Invalid patient Id")
		ac.logError(m, "Params", []zapcore.Field{zap.Any("ID", err.Message)})
		return response.RespondError(c, err)
	}

	dType := c.Query("type")
	from := c.Query("from")
	to := c.Query("to")

	ac.logInfo(m, "Input", []zapcore.Field{zap.Any("Id", patientId)})
	data := []models.Data{}
	if getErr := ac.UserService.GetSensorDataById(&data, patientId, dType, from, to); getErr.IsNotNull() {
		return response.RespondError(c, getErr)
	}
	return response.RespondSuccessJSON(c, &data)
}

func (tc *UserController) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(tc.logger, tc.name, method, message, keyValues)
}

func (tc *UserController) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(tc.logger, tc.name, method, message, keyValues)
}
