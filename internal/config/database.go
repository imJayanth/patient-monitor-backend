package config

import (
	"fmt"
	"log"
	"os"
	"patient-monitor-backend/internal/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func SetupDatabase() {
	appConfig := AppConfig

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True",
		appConfig.DBConfig.DBUSER,
		appConfig.DBConfig.DBPASSWORD,
		appConfig.DBConfig.DBHOST,
		appConfig.DBConfig.DBPORT,
		appConfig.DBConfig.DBNAME,
	)), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Println("Error Connecting to the DB", err.Error())
		panic(err)
	}

	DB = db
	sqldb, _ := db.DB()
	if err := sqldb.Ping(); err != nil {
		log.Fatal("Error while pinging the DB: ", err.Error())
	}
	if err := DB.Table(models.TABLE_USERS).AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error automigrating %s: %s", models.TABLE_USERS, err.Error())
	}
	if err := DB.Table(models.TABLE_CONTACTS).AutoMigrate(&models.Contact{}); err != nil {
		log.Fatalf("Error automigrating %s: %s", models.TABLE_CONTACTS, err.Error())
	}
	if err := DB.Table(models.TABLE_DATA).AutoMigrate(&models.Data{}); err != nil {
		log.Fatalf("Error automigrating %s: %s", models.TABLE_DATA, err.Error())
	}

	log.Printf("Connected to the DB: %s\n", appConfig.DBConfig.DBNAME)
}
