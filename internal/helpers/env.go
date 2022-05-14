package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(key string, defaultval string) string {
	if value, exists := os.LookupEnv(key); exists {
		fmt.Printf("Got value: %s for the key: %s \n", key, value)
		return value
	}
	return defaultval
}

func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, strconv.Itoa(defaultVal))
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
