package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"patient-monitor-backend/internal/errors"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CheckStringInList(s string, list []string) error {
	for _, v := range list {
		if strings.EqualFold(s, v) {
			return nil
		}
	}
	return fmt.Errorf("list doesn't contain the string")
}

func ForceCorrectTimeBeforeSave(t time.Time) time.Time {
	return t.Add(time.Hour * 5).Add(time.Minute * 30)
}

func ForceCorrectTimeAfterSave(t time.Time) time.Time {
	return t.Add(-time.Hour * 5).Add(-time.Minute * 30)
}

func ForceCorrectTimeAfterFind(timezone string, t time.Time) time.Time {
	loc, _ := time.LoadLocation(timezone)
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

func IsValidParametersByType(field interface{}) bool {
	value := reflect.ValueOf(field)
	valid := false

	switch value.Kind() {
	case reflect.String:
		if value.String() != "" {
			valid = true
		}
	case reflect.Int:
		if value.Int() > 0 {
			valid = true
		}
	case reflect.Float64:
		if value.Float() > 0 {
			valid = true
		}
	case reflect.Struct:
		if value.Field(0).Float() < -90 || value.Field(0).Float() > 90 || value.Field(0).Float() == 0 || value.Field(1).Float() < -180 || value.Field(1).Float() > 180 || value.Field(1).Float() == 0 {
			valid = false
		} else {
			valid = true
		}
		// can add many more fields
	}

	return valid
}

func LogInfo(logger *zap.Logger, calledFrom string, method string, message string, keyValues []zapcore.Field) {
	msg := calledFrom + ": " + method + ": " + message
	if keyValues == nil {
		logger.Info(msg)
	} else {
		logger.Info(msg, keyValues...)
	}
}

func LogError(logger *zap.Logger, calledFrom string, method string, message string, keyValues []zapcore.Field) {
	msg := calledFrom + ": " + method + ": " + message
	if keyValues == nil {
		logger.Error(msg)
	} else {
		logger.Error(msg, keyValues...)
	}
}

func ToJson(o interface{}) string {
	js, serr := json.Marshal(o)
	if serr != nil {
		log.Println("Error while marshalling : ", serr)
	}
	return string(js)
}

func ToJsonRaw(o interface{}) json.RawMessage {
	js, mErr := json.Marshal(o)
	if mErr != nil {
		log.Println("Error while marshalling Test: ", mErr)
	}
	return js
}

func ObjectToJson(x interface{}) string {
	js, err := json.Marshal(x)
	if err != nil {
		log.Println("Error while marshalling: ", err)
	}
	return string(js)
}

func MapErrorFromGorm(err error, msg ...string) errors.RestAPIError {

	if strings.Contains(err.Error(), "record not found") {
		return errors.NewNotFoundError("Invalid ID")
	} else if strings.Contains(err.Error(), "Duplicate entry") {
		return errors.NewDuplicateRecord("Duplicate Record")
	} else if strings.Contains(err.Error(), "Error 1452:") {
		return errors.NewInternalServerError(msg[0])
	}

	return errors.NewInternalServerError("something unexpected happen")

}
