package client

import (
	"io/ioutil"
	"math/rand"
	"netBlast/pkg"
	"netBlast/tools/scrapper"
	"strings"
)

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
