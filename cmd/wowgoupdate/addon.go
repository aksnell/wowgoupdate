package main

import (
	"errors"
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

	logger *log
	err    []error
}

func buildAddon(path string, errLogger *log) *addon {
	if _, err := os.Stat(getRelPath(path, "CHANGES.txt")); os.IsNotExist(err) {
		return nil
	}
	addonName := getAddonName(path)
	if addonName == "" {
		return nil
	}
	addon := &addon{
		Name:   addonName,
		Path:   path,
		logger: errLogger,
		err:    make([]error, 0),
	}
	addon.setLocalVersion()
	addon.setURL()
	addon.setOnlineVersion()
	return addon
}

func getAddonName(path string) string {
	changesPath := makeFilePath(path, `.toc`)
	if changesPath == "" {
		return ""
	}
	name, err := walkFile(changesPath, scanToc)
	if err != nil {
		return ""
	}
	return name
}

func (a *addon) setURL() {
	url := makeCurseURL(a.Name)
	if _, err := getURLBody(url); err == nil {
		a.URL = url
		return
	}
	a.log(a.Name, errors.New("failed to generate from Name, will try path"), 1)
	url = makeCurseURL(a.Path)
	if _, err := getURLBody(url); err == nil {
		a.URL = url
		return
	}
	a.log(a.Name, errors.New("failed to generate URL from name or path"), 0)
}

func (a *addon) setOnlineVersion() {
	body, err := getURLBody(a.URL + `\changes`)
	if err != nil {
		a.log(a.Name, err, 0)
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		a.log(a.Name, err, 0)
	}
	changeLine, err := doc.Find(".project-content.mg-t-1.pd-2 p").Eq(0).Html()
	if err != nil {
		a.log(a.Name, err, 0)
	}
	strippedLine := reAscii.ReplaceAllString(changeLine, " ")
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
		a.log(a.Name, err, 0)
	}
}

func (a *addon) log(prefix string, base error, level int) {
	a.logger.add(prefix, base, level)
}
