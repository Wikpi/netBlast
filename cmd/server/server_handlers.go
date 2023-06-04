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
	"strconv"
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
	users := []pkg.User{}

	userCnt, err := strconv.Atoi(database.QueryDB(s.db, "SELECT MAX(id) FROM users;"))
	if err != nil {
		panic(err.Error())
	}

	for x := 0; x < userCnt; x++ {
		users = append(users, pkg.User{
			Id:        x + 1,
			Name:      database.FindDBUserInfo(s.db, "name", "id", x+1),
			UserColor: database.FindDBUserInfo(s.db, "color", "id", x+1),
			Status:    database.FindDBUserInfo(s.db, "status", "id", x+1),
		})
	}
	userList := pkg.ParseToJson(users, "Couldnt parse to json.")

	w.Write(userList)
}

// Sends direct message to the recipient and the sender
func (s *serverInfo) directMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	pkg.HandleError(pkg.SvRegister+pkg.BadRead, err, 1)

	message := pkg.Message{}

	pkg.ParseFromJson(body, &message, "bad parse server")

	message.ReceiverColor = database.FindDBUserInfo(s.db, "color", "name", message.Receiver)

	var recConn, userConn *websocket.Conn

	for conn, value := range s.connections {
		if value == message.Receiver {
			recConn = conn
		}
		if value == message.Username {
			userConn = conn
		}
	}

	s.lock.RLock()
	pkg.WsWrite(recConn, message, "")
	s.lock.RUnlock()
	s.lock.RLock()
	pkg.WsWrite(userConn, message, "")
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
			fmt.Println("User left the server: ", s.connections[c])

			s.lock.Lock()
			database.UpdateDBUserInfo(s.db, "status", "name", "Offline", s.connections[c])
			s.lock.Unlock()
			return
		}

		if _, ok := s.connections[c]; !ok {
			s.lock.Lock()
			s.connections[c] = message.Username
			database.UpdateDBUserInfo(s.db, "status", "name", "Online", s.connections[c])
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
	for c, _ := range s.connections {
		pkg.WsWrite(c, message, pkg.SvMessage+pkg.BadWrite)
	}
	s.lock.RUnlock()
}
