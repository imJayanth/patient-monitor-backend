package services

import (
	"encoding/json"
	"fmt"
	"log"
	"patient-monitor-backend/internal/config"
	"patient-monitor-backend/internal/errors"
	"patient-monitor-backend/internal/helpers"
	"patient-monitor-backend/internal/models"
	"patient-monitor-backend/internal/repo"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type UserService struct {
	UserRepo *repo.UserRepo
	appConf  *config.AppConfigModel
	name     string
	logger   *zap.Logger
}

func NewUserService() *UserService {
	return &UserService{UserRepo: repo.NewUserRepo(), appConf: config.AppConfig, name: "UserService", logger: config.Logger}
}

func (ts *UserService) Login(user *models.User) errors.RestAPIError {
	m := "Login"
	ts.logInfo(m, "Start", nil)
	defer ts.logInfo(m, "Complete", nil)

	dbUser := models.User{}
	if err := ts.UserRepo.GetUserById(&dbUser, user.PatientId); err.IsNotNull() {
		return err
	}
	if dbUser.Password != user.Password {
		return errors.NewUnAuthorizedError("Invalid password")
	}
	return errors.NO_ERROR()
}

func (ts *UserService) AddContact(contact *models.Contact) errors.RestAPIError {
	m := "AddContact"
	ts.logInfo(m, "Start", nil)
	defer ts.logInfo(m, "Complete", nil)

	contacts := []models.Contact{}
	if err := ts.UserRepo.GetContactsById(&contacts, contact.PatientId); err.IsNotNull() {
		return err
	}
	for _, c := range contacts {
		if c.Type == contact.Type && c.Value == contact.Value {
			*contact = c
			return errors.NO_ERROR()
		}
	}
	return ts.UserRepo.SaveContact(contact)
}

func (ts *UserService) GetContactsById(contacts *[]models.Contact, patientId int) errors.RestAPIError {
	m := "GetContactsById"
	ts.logInfo(m, "Start", nil)
	defer ts.logInfo(m, "Complete", nil)

	if err := ts.UserRepo.GetContactsById(contacts, patientId); err.IsNotNull() {
		return err
	}
	return errors.NO_ERROR()
}

func (ts *UserService) GetSensorDataById(data *[]models.Data, patientId int, dType, from, to string) errors.RestAPIError {
	m := "GetSensorDataById"
	ts.logInfo(m, "Start", nil)
	defer ts.logInfo(m, "Complete", nil)
	sql := "select * from data where patient_id=" + fmt.Sprintf("%v", patientId)
	if dType != "" {
		sql += fmt.Sprintf(` AND type="%v"`, dType)
	}
	if from != "" && to != "" {
		// layout := "2006-01-02T15:04:05.000Z"
		// start, _ := time.Parse(layout, from)
		// end, _ := time.Parse(layout, to)
		sql += fmt.Sprintf(` AND timestamp between "%v" AND "%v"`, from, to)
	}
	log.Println(sql)
	dataMap := []map[string]interface{}{}
	if err := ts.UserRepo.ExecuteQuery(sql, &dataMap); err.IsNotNull() {
		return err
	}
	jsonString, _ := json.Marshal(dataMap)
	if err := json.Unmarshal(jsonString, &data); err != nil {
		return errors.NewInternalServerError("Error" + err.Error())
	}
	return errors.NO_ERROR()
}

func (ts *UserService) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(ts.logger, ts.name, method, message, keyValues)
}

func (ts *UserService) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(ts.logger, ts.name, method, message, keyValues)
}
