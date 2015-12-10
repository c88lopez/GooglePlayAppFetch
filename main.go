package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Apps list
type Apps struct {
	Apps []App
}

// App structure
type App struct {
	Name         string
	Developer    string
	BundleID     string
	AppStoreLink string
}

func main() {
	// Configs
	googlePlayURL := "https://play.google.com/store/"
	querySearch := os.Args[1]

	// Fetching the site
	doc, err := goquery.NewDocument(googlePlayURL + "search?q=" + url.QueryEscape(querySearch) + "&c=apps")
	if err != nil {
		log.Fatal(err)
	}

	appList := []App{}
	doc.Find(".card-list .card .card-content .details").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a.card-click-target").Attr("href")
		aux, _ := url.Parse(href)

		appID := strings.Trim(aux.Query()["id"][0], " ")
		appName := strings.Trim(s.Find("a.title").Text(), " ")
		appDeveloper := strings.Trim(s.Find(".subtitle-container .subtitle").Text(), " ")
		appLink := googlePlayURL + "apps/details?id=" + appID

		app := App{appName, appDeveloper, appID, appLink}
		appList = append(appList, app)
	})

	Apps := Apps{appList}
	js, err := json.Marshal(Apps)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", js)
}
