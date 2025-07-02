package utils

import "strconv"

func GetIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return num
}

func GetStringOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func GetBoolOrDefault(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	if value == "true" || value == "1" {
		return true
	}
	return false
}
