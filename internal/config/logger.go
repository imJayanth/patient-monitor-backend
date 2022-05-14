package config

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func SetupLogger() {
	appConfig := AppConfig

	writerSyncer := getLogWriterFile()
	encoder := getEncoder(appConfig)

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)

	Logger = zap.New(core, zap.AddCaller())
	// SugarredLogger = logger.Sugar()
}

func getEncoder(appConf *AppConfigModel) zapcore.Encoder {
	if appConf.Environment == "development" {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.FunctionKey = ""
		encoderConfig.CallerKey = ""
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriterFile() zapcore.WriteSyncer {
	file, err := os.OpenFile("./patient-monitor-backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return zapcore.AddSync(file)
}
