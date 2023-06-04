package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"netBlast/pkg"
	"netBlast/tools/scrapper"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Pings server
func pingServer() int {
	res := handleGetRequest("http://"+pkg.ServerURL+"/", "Ping: ")

	return res.StatusCode
}

// Receives user list from server
func getUserList(m *model) {
	res := handleGetRequest("http://"+pkg.ServerURL+"/userList", pkg.ClList)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(pkg.Cl + pkg.BadRead)
	}

	m.lock.Lock()
	pkg.ParseFromJson(resBody, &m.userList.users, pkg.ClList+pkg.BadParseFrom)
	m.lock.Unlock()

	res.Body.Close()
}

// Handles POST request to server
func handlePostRequest(data []byte, URL string, incomingErr string) *http.Response {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(data))
	if err != nil {
		pkg.LogError(err)
		fmt.Println(incomingErr + pkg.BadReq)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(incomingErr + pkg.BadRes)
	}

	return res
}

// Handles Get request to server
func handleGetRequest(URL string, incomingErr string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(incomingErr + pkg.BadReq)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(incomingErr + pkg.BadRes)
	}

	return res
}

// Change screen logic
func changeScreen(m *model, nextScreen string) (*model, tea.Cmd) {
	if m.screen == "register" {
		return m, nil
	}

	if m.screen != nextScreen {
		m.prevScreen = m.screen
		m.screen = nextScreen
		return m, nil
	}

	if nextScreen != "help" {
		m.screen = "chat"
		return m, nil
	}
	m.screen = m.prevScreen
	return m, nil
}

// Picks one random color from the scrapped color list
func getColor() string {
	path := pkg.Scrapper + "/colors.txt"

	body, err := ioutil.ReadFile(path)
	if err != nil {
		pkg.LogError(err)
		fmt.Println(pkg.Cl + pkg.BadRead)
	}
	if body == nil {
		return "#FFF"
	}

	colors := strings.Split(string(body), ", ")

	return colors[rand.Intn(len(colors))]
}

// Scrapes colors using Autolycus module
func useAutolycus() {
	scrapper.Scrape()
}
