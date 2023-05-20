package client

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"netBlast/pkg"
	"netBlast/tools/scrapper"
	"strings"
)

// Pings server
func pingServer() int {
	res := handleGetRequest("http://"+pkg.ServerURL+"/", "Ping: ")

	return res.StatusCode
}

// Receives user list from server
func getUserList(m *model) {
	res := handleGetRequest("http://"+pkg.ServerURL+"/userList", "Client/ListUsers: ")

	resBody, err := ioutil.ReadAll(res.Body)
	pkg.HandleError(pkg.Cl+pkg.BadRead, err)

	m.lock.Lock()
	pkg.ParseFromJson(resBody, &m.userList.users, pkg.Cl+pkg.BadParse)
	m.lock.Unlock()

	res.Body.Close()
}

// Handles POST request to server
func handlePostRequest(data []byte, URL string, incomingErr string) *http.Response {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))
	pkg.HandleError(incomingErr+pkg.BadReq, err, 0)

	res, err := http.DefaultClient.Do(req)
	pkg.HandleError(incomingErr+pkg.BadRes, err, 0)

	return res
}

// Handles Get request to server
func handleGetRequest(URL string, incomingErr string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	pkg.HandleError(incomingErr+pkg.BadReq, err, 0)

	res, err := http.DefaultClient.Do(req)
	pkg.HandleError(incomingErr+pkg.BadRes, err, 0)

	return res
}

// Picks one random color from the scrapped color list
func getColor() string {
	path := pkg.Scrapper + "/colors.txt"

	body, err := ioutil.ReadFile(path)
	pkg.HandleError(pkg.Cl+pkg.BadOpen+": "+path, err, 1)

	colors := strings.Split(string(body), ", ")

	return colors[rand.Intn(len(colors))]
}

// Scrapes colors using Autolycus module
func useAutolycus() {
	scrapper.Scrape()
}
