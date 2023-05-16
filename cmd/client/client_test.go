package main

import (
	"io/ioutil"
	"netBlast/pkg"
	"netBlast/tools/scrapper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetColor(t *testing.T) {
	testColor := getColor()

	assert.NotEmpty(t, testColor)
}

func Test_UseAutolycus(t *testing.T) {
	scrapper.Scrape()

	body, err := ioutil.ReadFile(pkg.Scrapper + "/colors.txt")
	pkg.HandleError(pkg.Cl+pkg.BadOpen, err, 0)

	assert.NotEmpty(t, body)
}
