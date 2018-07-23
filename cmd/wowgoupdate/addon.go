package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type addon struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	URL     string `json:"url"`
	Version string `json:"version"`
	Latest  string `json:"latest"`
	Current bool   `json:"current"`

	err []error
}

func buildAddon(path string) (*addon, error) {
	if _, err := os.Stat(getRelPath(path, "CHANGES.txt")); os.IsNotExist(err) {
		return nil, err
	}
	addonName, err := getAddonName(path)
	if addonName == "" {
		return nil, err
	}
	addon := &addon{
		Name: addonName,
		Path: path,
	}
	addon.setLocalVersion()
	addon.setURL()
	addon.setOnlineVersion()
	return addon, nil
}

func getAddonName(path string) (string, error) {
	changesPath := makeFilePath(path, `.toc`)
	if changesPath == "" {
		return "", errors.New("Could not make .toc filepath from," + path)
	}
	name, err := walkFile(changesPath, scanToc)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (a *addon) setURL() {
	url := makeCurseURL(a.Name)
	if _, err := getURLBody(url); err == nil {
		a.URL = url
		return
	}
	url = makeCurseURL(a.Path)
	if _, err := getURLBody(url); err == nil {
		a.URL = url
	}
}

func (a *addon) setOnlineVersion() {
	fmt.Println(a.URL)
	body, err := getURLBody(a.URL + `\changes`)
	if err != nil {
		fmt.Println(err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	changeLine, err := doc.Find(".project-content.mg-t-1.pd-2 p").Eq(0).Html()
	if err != nil {
		fmt.Println(err)
		return
	}
	strippedLine := reASCII.ReplaceAllString(changeLine, " ")
	splitLines := reBrTag.Split(strippedLine, 3)
	for i := range splitLines {
		if reAlphaNum.MatchString(splitLines[i]) {
			a.Latest = splitLines[i]
			a.Current = a.Latest == a.Version
			break
		}
	}
}

func (a *addon) setLocalVersion() {
	var err error
	a.Version, err = walkFile(getRelPath(a.Path, "CHANGES.txt"), scanChanges)
	if err != nil {
		fmt.Println(err)
	}
}
