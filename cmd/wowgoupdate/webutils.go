package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDocumentFromURL(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

//isValidURL returns true if an HTTP response returns status code 200 and no errors.
func isValidURL(URL string) bool {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		return false
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	return resp.StatusCode == 200
}

//return formatted Curse URL from Addon.Name or Addon.Path
func makeCurseURL(nameOrPath string) string {
	formattedBaseURL := strings.Join(reAlphaNum.FindAllString(filepath.Base(nameOrPath), -1), "-")
	return strings.Join([]string{curseURL, formattedBaseURL}, "")
}
