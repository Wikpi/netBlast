package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"netBlast/pkg"
	"netBlast/tools/database"
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
		//s.users = append(s.users, client)
		database.InsertDBUser(s.db, client)
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

	s.readMessage(c)
}

// Sends back the list of users
func (s *serverInfo) sendUserList(w http.ResponseWriter, r *http.Request) {
	users := pkg.ParseToJson(s.users, "Server/SendUsers: couldnt parse to json.")

	w.Write(users)
}

// Sends direct message to the recipient and the sender
func (s *serverInfo) directMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	pkg.HandleError(pkg.SvRegister+pkg.BadRead, err, 1)

	message := pkg.Message{}

	pkg.ParseFromJson(body, &message, "bad parse server")

	message.ReceiverColor = s.users[findUser(message.Receiver, s)].UserColor

	s.lock.RLock()
	pkg.WsWrite(s.users[findUser(message.Receiver, s)].Conn, message, "")
	s.lock.RUnlock()
	s.lock.RLock()
	pkg.WsWrite(s.users[findUser(message.Username, s)].Conn, message, "")
	s.lock.RUnlock()
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
			name := database.FindDBUserInfo(s.db, "name", "conn", c)

			fmt.Println("User left the server: ", name)

			s.lock.Lock()
			database.UpdateDBUserInfo(s.db, name, "status", "Offline")
			s.lock.Unlock()
			//s.users = append(s.users[:userIdx], s.users[userIdx+1:]...)
			return
		}

		if !database.CheckDBUserConn(s.db, message.Username) {
			s.lock.Lock()
			database.UpdateDBUserInfo(s.db, message.Username, "conn", c)
			s.lock.Unlock()
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
