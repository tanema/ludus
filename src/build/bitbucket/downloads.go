package bitbucket

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

const downloadsURL = "https://api.bitbucket.org/2.0/repositories/rude/love/downloads"
const releaseURL = "https://love2d.org/releases.xml"

type downloadsPage struct {
	PageCount int        `json:"pagelen"`
	Page      int        `json:"page"`
	Next      string     `json:"next"`
	Downloads []download `json:"values"`
}

type download struct {
	Name  string `json:"name"`
	Links links  `json:"links"`
}

type links struct {
	Self link `json:"self"`
}

type link struct {
	Href string `json:"href"`
}

func getDownloadPage(page string) (downloadsPage, error) {
	resp, err := http.Get(page)
	if err != nil {
		return downloadsPage{}, err
	}
	defer resp.Body.Close()

	var downloads downloadsPage
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&downloads); err != nil {
		return downloadsPage{}, err
	}

	return downloads, nil
}

type rss struct {
	XMLName xml.Name  `xml:"rss"`
	Channel []channel `xml:"channel"`
}

type channel struct {
	Items []item `xml:"item"`
}

type item struct {
	Title string `xml:"title"`
}

func getRecentLoveVersion() (string, error) {
	resp, err := http.Get(releaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var updates rss
	decoder := xml.NewDecoder(resp.Body)
	if err = decoder.Decode(&updates); err != nil {
		return "", nil
	}

	fullVersion := updates.Channel[0].Items[0].Title
	parts := strings.Split(fullVersion, " ")
	return parts[0], nil
}
