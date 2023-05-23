package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"netBlast/pkg"
	"os"
	"os/signal"
	"time"

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

// Shutdowns the server
func (server *serverInfo) serverShutdown() {
	signal.Notify(server.shutdown, os.Interrupt)
	<-server.shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.s.Shutdown(ctx)
	pkg.HandleError(pkg.Sv+pkg.BadClose, err, 1)
}

// Reads received messages
func (s *serverInfo) readMessage(c *websocket.Conn) {
	for {
		message := pkg.WsRead(c, pkg.SvMessage+pkg.BadRead)
		if (message == pkg.Message{}) {
			if userIdx := findUser(c, s); userIdx != -1 {
				fmt.Println("User left the server: ", s.users[userIdx].Name)
				s.users = append(s.users[:userIdx], s.users[userIdx+1:]...)
				return
			}
		}

		s.lock.Lock()
		s.messages = append(s.messages, message)
		s.lock.Unlock()

		s.writeToAll(message)
	}
}

// Writes user message to all other connections
func (s *serverInfo) writeToAll(message pkg.Message) {
	s.lock.RLock()
	for _, ic := range s.users {
		pkg.WsWrite(ic.Conn, message, pkg.SvMessage+pkg.BadWrite)
	}
	s.lock.RUnlock()
}
