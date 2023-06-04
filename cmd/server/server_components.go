package server

import (
	"net/http"
	"netBlast/tools/database"
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

		if name := database.FindDBUserInfo(s.db, "name", "name", name); name == "" {
			errMsg = ""
			statusCode = http.StatusAccepted
		} else {
			errMsg = "Name already exists. "
			statusCode = http.StatusNotAcceptable
		}
	}

	return errMsg, statusCode
}
