package client

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"netBlast/cmd/server"
	"netBlast/pkg"
	"netBlast/tools/scrapper"

	"github.com/stretchr/testify/assert"
)

func Test_GetColor(t *testing.T) {
	testColor := getColor()

	assert.NotEmpty(t, testColor)
}

func Test_UseAutolycus(t *testing.T) {
	scrapper.Scrape()

	body, err := ioutil.ReadFile("../.." + pkg.Scrapper + "/colors.txt")
	pkg.HandleError(pkg.Cl+pkg.BadOpen, err, 0)

	assert.NotEmpty(t, body)
}
func Test_HandleHTTPRequest(t *testing.T) {
	const URL = "http://" + pkg.ServerURL + "/register"

	shutdown := make(chan os.Signal)

	go func() {
		server.Server(shutdown)
	}()

	name := pkg.Name{Name: "Bobby"}
	data := pkg.ParseToJson(name, "test")

	res := handlePostRequest(data, URL, "")
	assert.NotEmpty(t, res)

	shutdown <- syscall.SIGINT
}
