package env

import (
	"os"
	"strconv"
)

func getFromEnv(variableName string) string {
	return os.Getenv(variableName)
}

func String(variableName string, defaultValue string) string {
	strValue := getFromEnv(variableName)
	if strValue == "" {
		strValue = defaultValue
	}
	return strValue
}

func Int(variableName string, defaultValue int) int {
	strValue := getFromEnv(variableName)
	value := defaultValue
	if strValue != "" {
		newVal, err := strconv.Atoi(strValue)
		if err == nil {
			value = newVal
		}
	}
	return value
}
