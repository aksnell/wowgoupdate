package main

import (
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
		return nil // path/to/whatever does not exist
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
	changesPath := getRelPath(path, makeFilePath(path, `.toc`))
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
	_, err := getURLBody(url)
	if err != nil {
		url = makeCurseURL(a.Path)
		_, err = getURLBody(url)
		if err == nil {
			a.URL = url
		} else {
			a.URL = "URL NOT FOUND."
		}
	} else {
		a.URL = url
	}
}

func (a *addon) setOnlineVersion() {
	body, err := getURLBody(a.URL + `\changes`)
	if err != nil {
		a.log("setOnlineVersion", err, 0)
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		a.log("setOnlineVersion", err, 0)
	}
	changeLine, err := doc.Find(".project-content.mg-t-1.pd-2 p").Eq(0).Html()
	if err != nil {
		a.log("setOnlineVersion", err, 0)
	}
	//changeLine = regexp.MustCompile(`[^[:ascii:]]`).ReplaceAllString(changeLine, " ")
	splitLines := reBrTag.Split(changeLine, 3)
	for i := range splitLines {
		if reAlphaNum.MatchString(splitLines[i]) {
			a.Current = a.Latest == a.Version
			a.Latest = splitLines[i]
			break
		}
	}
}

func (a *addon) setLocalVersion() {
	var err error
	a.Version, err = walkFile(getRelPath(a.Path, "CHANGES.txt"), scanChanges)
	if err != nil {
		a.log("setLocalVersion", err, 0)
	}
}

func (a *addon) log(prefix string, base error, level int) {
	a.logger.add(prefix, base, level)
}
