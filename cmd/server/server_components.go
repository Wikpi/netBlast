package server

import (
	"net/http"
	"unicode/utf8"
)

// Validates username
func checkName(name string, s *serverInfo) (string, int) {
	errMsg := ""
	statusCode := 0

	if utf8.RuneCountInString(name) < 3 {
		errMsg = "Name too short. "
		statusCode = http.StatusNotAcceptable
	} else if utf8.RuneCountInString(name) > 10 {
		errMsg = "Name too long. "
		statusCode = http.StatusNotAcceptable
	} else {
		if user := findUser(name, s); user == -1 {
			errMsg = ""
			statusCode = http.StatusAccepted
		} else {
			errMsg = "Name already exists. "
			statusCode = http.StatusNotAcceptable
		}
	}

	return errMsg, statusCode
}

// Finds name in user slice
func findUser(key interface{}, s *serverInfo) int {
	for idx, user := range s.users {
		if user.Name == key || user.Conn == key {
			return idx
		}
	}
	return -1
}
