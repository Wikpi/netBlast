package server

import (
	"net/http"
	"netBlast/pkg"
	"os"
	"sync"
)

var (
	pingHandler     = "/"
	registerHandler = "/register"
	sessionHandler  = "/message"
	userListHandler = "/userList"
	dmHandler       = "/dmUser"
)

// Server struct that holds all the essential info
type serverInfo struct {
	s        http.Server
	messages []pkg.Message
	users    []pkg.User
	lock     sync.RWMutex
	mux      *http.ServeMux
	shutdown chan os.Signal
}
