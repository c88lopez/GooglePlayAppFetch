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
	Name         string `json:"name"`
	Developer    string `json:"developer"`
	BundleID     string `json:"bundleID"`
	AppStoreLink string `json:"appStoreLink"`
	Icon         string `json:"icon"`
}

func main() {
	// Configs
	googlePlayURL := "https://play.google.com/store/"

	if len(os.Args) != 2 {
		fmt.Println("I just need a string to search...")
		os.Exit(1)
	}

	querySearch := os.Args[1]

	// Fetching the site
	doc, err := goquery.NewDocument(googlePlayURL + "search?q=" + url.QueryEscape(querySearch) + "&c=apps")
	if err != nil {
		log.Fatal(err)
	}

	appList := []App{}
	doc.Find(".card-list .card .card-content").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a.card-click-target").Attr("href")
		aux, _ := url.Parse(href)

		appID := strings.Trim(aux.Query()["id"][0], " ")
		appName := strings.Trim(s.Find("a.title").Text(), " ")
		appDeveloper := strings.Trim(s.Find(".subtitle-container .subtitle").Text(), " ")
		appLink := googlePlayURL + "apps/details?id=" + appID
		appIconRaw, _ := s.Find("img.cover-image").Attr("src")

		appIcon, _ := url.Parse(appIconRaw)

		if "" == appIcon.Scheme {
			appIcon.Scheme = "http"
		}

		app := App{appName, appDeveloper, appID, appLink, appIcon.String()}
		appList = append(appList, app)
	})

	Apps := Apps{appList}
	js, err := json.Marshal(Apps)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", js)
}
