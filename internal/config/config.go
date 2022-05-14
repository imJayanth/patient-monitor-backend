package config

import (
	"os"
	"patient-monitor-backend/internal/helpers"
	"time"
)

type dbConfig struct {
	DBUSER     string
	DBPASSWORD string
	DBHOST     string
	DBNAME     string
	DBPORT     string
}

type serverConfig struct {
	APIPORT  int
	APINAME  string
	TIMEZONE string
}

type authSecrets struct {
	JWTSECRETKEY string
}

type AppConfigModel struct {
	DBConfig        dbConfig
	AuthSecrets     authSecrets
	ServerConfig    serverConfig
	Environment     string
	ApplicationPath string
}

var AppConfig *AppConfigModel

func SetupConfig() {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	AppConfig = &AppConfigModel{
		DBConfig: dbConfig{
			DBUSER:     helpers.GetEnv("DBUSER", "root"),
			DBPASSWORD: helpers.GetEnv("DBPASSWORD", "abcd1234"),
			DBHOST:     helpers.GetEnv("DBHOST", "localhost"),
			DBNAME:     helpers.GetEnv("DBNAME", "patient-monitor"),
			DBPORT:     helpers.GetEnv("DBPORT", "3306"),
		},
		ServerConfig: serverConfig{
			APIPORT:  helpers.GetEnvAsInt("APIPORT", 8081),
			APINAME:  helpers.GetEnv("APINAME", "patient-monitor-backend"),
			TIMEZONE: helpers.GetEnv("TIMEZONE", "Asia/Kolkata"),
		},
		AuthSecrets: authSecrets{
			JWTSECRETKEY: helpers.GetEnv("JWTSECRETKEY", ""),
		},
		Environment:     "development",
		ApplicationPath: helpers.GetEnv("ApplicationPath", currentPath),
	}
}

//Time zone related time
var loc *time.Location

func (s *AppConfigModel) CurrentTime() time.Time {
	if loc == nil {
		location, err := time.LoadLocation(s.ServerConfig.TIMEZONE)
		if err != nil {
			panic(err)
		}
		loc = location
	}
	t := time.Now()
	return t.In(loc)
}

func (s *AppConfigModel) Timezone() *time.Location {
	if loc == nil {
		location, err := time.LoadLocation(s.ServerConfig.TIMEZONE)
		if err != nil {
			panic(err)
		}
		loc = location
	}
	return loc
}
