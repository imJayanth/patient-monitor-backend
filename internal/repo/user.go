package repo

import (
	"fmt"
	"patient-monitor-backend/internal/config"
	"patient-monitor-backend/internal/errors"
	"patient-monitor-backend/internal/helpers"
	"patient-monitor-backend/internal/models"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB      *gorm.DB
	appConf *config.AppConfigModel
	name    string
	logger  *zap.Logger
}

func NewUserRepo() *UserRepo {
	return &UserRepo{DB: config.DB, appConf: config.AppConfig, name: "UserRepo", logger: config.Logger}
}

func (tr *UserRepo) SaveContact(contact *models.Contact) errors.RestAPIError {
	m := "SaveContact"

	if err := tr.DB.Table(models.TABLE_CONTACTS).Save(contact).Error; err != nil {
		saveErr := errors.NewInternalServerError(fmt.Sprintf("Could not save contact: %+v", err.Error()))
		tr.logError(m, "SaveErr", []zapcore.Field{zap.Any("contact", saveErr.Message)})
		return saveErr
	}
	tr.logInfo(m, "Saved contact", []zapcore.Field{zap.Any("id", contact.ToJsonRaw())})
	return errors.NO_ERROR()
}

func (tr *UserRepo) GetUserById(user *models.User, id int) errors.RestAPIError {
	m := "GetUserById"
	if err := tr.DB.Table(models.TABLE_USERS).Where(&models.User{PatientId: user.PatientId}).Find(&user).Error; err != nil {
		getErr := errors.NewNotFoundError(fmt.Sprintf("Could not find User- %v", err.Error()))
		tr.logError(m, "GetErr", []zapcore.Field{zap.Any("User", getErr.Message)})
		return getErr
	} else if user.IsNull() {
		getErr := errors.NewNotFoundError(fmt.Sprintf("Could not find User"))
		tr.logError(m, "GetErr", []zapcore.Field{zap.Any("User", getErr.Message)})
		return getErr
	}
	tr.logInfo(m, "From DB", []zapcore.Field{zap.Any("User", user.ToJsonRaw())})
	return errors.NO_ERROR()
}

func (tr *UserRepo) GetContactsById(contacts *[]models.Contact, patientId int) errors.RestAPIError {
	m := "GetContactsById"
	if err := tr.DB.Table(models.TABLE_CONTACTS).Where(&models.Contact{PatientId: patientId}).Find(&contacts).Error; err != nil {
		getErr := errors.NewNotFoundError(fmt.Sprintf("Could not find contacts- %v", err.Error()))
		tr.logError(m, "GetErr", []zapcore.Field{zap.Any("User", getErr.Message)})
		return getErr
	} else if len(*contacts) == 0 {
		getErr := errors.NewNotFoundError("Couldn't find any contacts")
		tr.logError(m, "GetErr", []zapcore.Field{zap.Any("contacts", getErr.Message)})
		return getErr
	}
	tr.logInfo(m, "From DB", []zapcore.Field{zap.Any("contacts", helpers.ToJsonRaw(contacts))})
	return errors.NO_ERROR()
}

func (o *UserRepo) ExecuteQuery(sql string, r *[]map[string]interface{}) errors.RestAPIError {
	// var r []map[string]interface{}
	// sqldb :=
	// rows, err := o.DB.CommonDB().Query(sql)
	dbSql, e := o.DB.DB()
	if e != nil {
		return errors.NewInternalServerError("Error executing sql e- " + e.Error())
	}
	rows, err := dbSql.Query(sql)
	if err != nil {
		return errors.NewInternalServerError("Error executing sql - " + err.Error())
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	cols, _ := rows.Columns()
	colTypes, _ := rows.ColumnTypes()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		var columns []interface{}
		for _, ct := range colTypes {
			if ct.DatabaseTypeName() == "VARCHAR" {
				columns = append(columns, new(string))
			} else if ct.DatabaseTypeName() == "TIMESTAMP" {
				columns = append(columns, time.Time{})
			} else if ct.DatabaseTypeName() == "DATETIME" {
				columns = append(columns, time.Time{})
			} else if ct.DatabaseTypeName() == "DECIMAL" {
				columns = append(columns, 0.0)
			} else if ct.DatabaseTypeName() == "DOUBLE" {
				columns = append(columns, 0.0)
			} else if ct.DatabaseTypeName() == "BIGINT" {
				columns = append(columns, 0)
			} else if ct.DatabaseTypeName() == "TINYINT" {
				columns = append(columns, new(bool))
			} else if ct.DatabaseTypeName() == "INT" {
				columns = append(columns, 0)
			} else {
				columns = append(columns, new(interface{}))
			}
		}

		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return errors.NewInternalServerError("Error scanning rows - " + err.Error())
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			v := *val
			switch v.(type) {
			case []uint8:
				m[colName] = string(v.([]byte))
			default:
				m[colName] = *(columnPointers[i].(*interface{}))
			}
			v2 := m[colName]
			ct := colTypes[i]
			switch v2.(type) {
			case string:
				s := v2.(string)
				if ct.DatabaseTypeName() == "VARCHAR" {
					m[colName] = s
				} else if ct.DatabaseTypeName() == "TIMESTAMP" {
					layout := "2006-01-02T15:04:05.000Z"
					t, _ := time.Parse(layout, s)
					m[colName] = helpers.ForceCorrectTimeAfterFind(o.appConf.ServerConfig.TIMEZONE, t)
				} else if ct.DatabaseTypeName() == "DATETIME" {
					layout := "2006-01-02T15:04:05.000Z"
					t, _ := time.Parse(layout, s)
					m[colName] = helpers.ForceCorrectTimeAfterFind(o.appConf.ServerConfig.TIMEZONE, t)
				} else if ct.DatabaseTypeName() == "DECIMAL" {
					f, _ := strconv.ParseFloat(s, 64)
					m[colName] = f
				} else if ct.DatabaseTypeName() == "DOUBLE" {
					f, _ := strconv.ParseFloat(s, 64)
					m[colName] = f
				} else if ct.DatabaseTypeName() == "BIGINT" {
					n, _ := strconv.ParseInt(s, 10, 64)
					m[colName] = n
				} else if ct.DatabaseTypeName() == "TINYINT" {
					n, _ := strconv.ParseBool(s)
					m[colName] = n
				} else if ct.DatabaseTypeName() == "INT" {
					n, _ := strconv.ParseInt(s, 10, 64)
					m[colName] = n
				} else {
					m[colName] = s
				}
			}
		}
		*r = append(*r, m)
	}
	if len(*r) == 0 {
		return errors.NewNotFoundError("Could not find data")
	}
	return errors.NO_ERROR()
}

func (tr *UserRepo) ForceCorrectTimeBeforeSave(data *models.Data) {
	data.Timestamp = helpers.ForceCorrectTimeBeforeSave(data.Timestamp)
}

func (tr *UserRepo) ForceCorrectTimeAfterSave(data *models.Data) {
	data.Timestamp = helpers.ForceCorrectTimeAfterSave(data.Timestamp)
}

func (tr *UserRepo) ForceCorrectTimeAfterFind(data *models.Data) {
	data.Timestamp = helpers.ForceCorrectTimeAfterFind(tr.appConf.ServerConfig.TIMEZONE, data.Timestamp)
}

func (tr *UserRepo) logInfo(method, message string, keyValues []zapcore.Field) {
	helpers.LogInfo(tr.logger, tr.name, method, message, keyValues)
}

func (tr *UserRepo) logError(method, message string, keyValues []zapcore.Field) {
	helpers.LogError(tr.logger, tr.name, method, message, keyValues)
}
