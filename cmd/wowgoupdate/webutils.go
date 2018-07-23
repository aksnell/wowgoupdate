package main

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getUpdateURL(url string) string {
	body, err := getURLBody(url + `\files`)
	if err != nil {
		return "error at get URLBody"
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "err at make reader"
	}
	updateURL, _ := doc.Find(".button.button--download.download-button.mg-r-05").Attr("href")
	if err != nil {
		return "err at find"
	}
	return updateURL
}

func makeCurseURL(path string) string {
	base := filepath.Base(path)
	baseURL := strings.Join(reAlphaNum.FindAllString(base, -1), "-")
	return strings.Join([]string{curseURL, baseURL}, "")
}

func getURLBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("getURLBody:bad response code " + string(resp.StatusCode))
	}
	return resp.Body, nil
}
