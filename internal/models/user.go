package models

import (
	"encoding/json"
	"log"
	"reflect"
	"time"
)

const TABLE_USERS = "users"
const TABLE_CONTACTS = "contacts"
const TABLE_DATA = "data"

func (User) TableName() string {
	return TABLE_USERS
}

func (Contact) TableName() string {
	return TABLE_CONTACTS
}

func (Data) TableName() string {
	return TABLE_DATA
}

type User struct {
	PatientId int    `json:"patient_id" gorm:"column:patient_id;primaryKey;autoIncrement"`
	Password  string `json:"password"`
}

type Contact struct {
	Id        int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PatientId int    `json:"patient_id"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type Data struct {
	Id        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	PatientId int       `json:"patient_id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Value     string    `json:"value"`
}

func (t *User) ToJson() string {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling User: ", mErr)
	}
	return string(js)
}

func (t *User) ToJsonRaw() json.RawMessage {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling User: ", mErr)
	}
	return js
}

func (t User) IsNull() bool {
	return reflect.DeepEqual(t, User{})
}

func (t *Contact) ToJson() string {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling Contact: ", mErr)
	}
	return string(js)
}

func (t *Contact) ToJsonRaw() json.RawMessage {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling Contact: ", mErr)
	}
	return js
}

func (t Contact) IsNull() bool {
	return reflect.DeepEqual(t, Contact{})
}

func (t *Data) ToJson() string {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling Data: ", mErr)
	}
	return string(js)
}

func (t *Data) ToJsonRaw() json.RawMessage {
	js, mErr := json.Marshal(t)
	if mErr != nil {
		log.Println("Error while marshalling Data: ", mErr)
	}
	return js
}

func (t Data) IsNull() bool {
	return reflect.DeepEqual(t, Data{})
}
