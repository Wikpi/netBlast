package scrapper

import autolycus "github.com/Wikpi/Autolycus/pkg"

func Scrape() {
	colors := []string{}

	// website to scrape
	url := "https://htmlcolorcodes.com/colors/"
	// Arguments to scrape (tag, key, value !)
	arg := []string{"td", "class", "color-table__cell--hex"}
	// Write path of the txt file
	path := "./tools/scrapper/colors.txt"

	// Iniates the scrapper i.e. get the html string and parses it
	doc := autolycus.Initiate(url)
	// Scrapes the data
	autolycus.Scrape(&colors, doc, arg)

	autolycus.WriteData(path, colors)
}
