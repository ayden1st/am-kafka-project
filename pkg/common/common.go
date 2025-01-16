package common

import (
	"fmt"
	"os"
	"strings"
)

// MapStringJson - function that converts a map of strings to a JSON string.
// It takes a map of strings as an argument and returns a JSON string.
func MapStringJson(array map[string]string) string {
	str := "{"
	for key, value := range array {
		value = strings.ReplaceAll(value, "\n", "")                            // Remove newline characters from the value
		str += fmt.Sprint("\"" + key + "\"" + ":" + "\"" + value + "\"" + ",") // Append the key-value pair to the JSON string
	}
	str = str[:len(str)-1] + "}" // Remove the trailing comma and close the JSON object
	return str                   // Return the JSON string
}

// GetEnv - function that retrieves the value of an environment variable.
// It takes two arguments: key (the name of the environment variable) and fallback (the default value to return if the environment variable is not set).
// The function returns the value of the environment variable or the fallback value if the environment variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value // If the environment variable is set, return its value
	}
	return fallback // If the environment variable is not set, return the fallback value
}
