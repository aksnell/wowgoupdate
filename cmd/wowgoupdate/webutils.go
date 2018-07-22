package main

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

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
