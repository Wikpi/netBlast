package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"netBlast/pkg"

	"nhooyr.io/websocket"
)

// Pings the server
func ping(w http.ResponseWriter, r *http.Request) {

	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	fmt.Println("Update call from: ", IPAddress)

	w.WriteHeader(http.StatusOK)
}

// Registers new users
func (s *serverInfo) registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	pkg.HandleError(pkg.SvRegister+pkg.BadRead, err, 1)

	name := pkg.Name{}

	pkg.ParseFromJson(body, &name, pkg.SvRegister+pkg.BadParse)

	s.lock.Lock()
	errMsg, status := checkName(name.Name, s)
	s.lock.Unlock()

	if status == http.StatusAccepted {
		fmt.Println("New registered user: ", name.Name)

		client := pkg.User{
			Name:      name.Name,
			UserColor: name.Color,
		}

		s.lock.Lock()
		s.users = append(s.users, client)
		s.lock.Unlock()
	}

	w.WriteHeader(status)
	if errMsg != "" {
		data := pkg.ParseToJson(errMsg, pkg.SvRegister+pkg.BadParse)
		w.Write(data)
	}
}

// Handle websocket connection
func (s *serverInfo) handleSession(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	pkg.HandleError(pkg.SvMessage+pkg.BadConn, err, 1)

	defer c.Close(websocket.StatusInternalError, "")

	s.lock.Lock()
	s.users[len(s.users)-1].Conn = c
	s.users[len(s.users)-1].Status = "Online"
	s.lock.Unlock()

	s.readMessage(c)
}

// Sends back the list of users
func (s *serverInfo) sendUserList(w http.ResponseWriter, r *http.Request) {
	users := pkg.ParseToJson(s.users, "Server/SendUsers: couldnt parse to json.")

	w.Write(users)
}
