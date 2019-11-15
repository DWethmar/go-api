package main

import (
	"strconv"
	"strings"
)

// URLPathPart gets the url path part by index.
func URLPathPart(path string, index int) (string, error) {
	values := strings.Split(path, "/")
	var parts []string
	for _, value := range values {
		if value != "" {
			parts = append(parts, value)
		}
	}
	return parts[index], nil
}

// URLPathPartInt gets the url path part by index as int.
func URLPathPartInt(urlStr string, index int) (int, error) {
	part, err := URLPathPart(urlStr, index)
	if err != nil {
		return 0, err
	}
	value, parseErr := strconv.Atoi(part)
	if parseErr != nil {
		return 0, parseErr
	}
	return value, err
}
