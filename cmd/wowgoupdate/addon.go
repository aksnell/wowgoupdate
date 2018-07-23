package main

import (
	"fmt"
)

type addon struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	URL     string `json:"url"`
	Version string `json:"version"`
	Latest  string `json:"latest"`
}

func buildAddon(path string) *addon {
	if isValid := fileExists(makeSpecificPath(path, `CHANGES.txt`)); isValid {
		return nil
	}
	fmt.Println(path)
	addon := &addon{
		Path: path,
	}
	addon.setName()
	addon.setLocalVersion()
	addon.setURL()
	addon.setOnlineVersion()
	return addon
}

func (a *addon) setName() {
	var err error
	a.Name, err = walkFile(makeGenericPath(a.Path, `.toc`), scanToc)
	if err != nil {
		fmt.Println(a.Path, "could not set name", err)
	}
}

func (a *addon) setURL() {
	urlFromName := makeCurseURL(a.Name)
	if isValid := isValidURL(urlFromName); isValid {
		a.URL = urlFromName
		return
	}
	urlFromPath := makeCurseURL(a.Path)
	if isValid := isValidURL(urlFromPath); isValid {
		a.URL = urlFromPath
		return
	}
	fmt.Println(a.Name, "could not set valid url")
}

func (a *addon) setOnlineVersion() {
	doc, err := getDocumentFromURL(a.URL + `/changes`)
	changesFile, err := doc.Find(".project-content.mg-t-1.pd-2 p").Eq(0).Html()
	if err != nil {
		fmt.Println(a.Name, "could not set online version", err)
	}
	splitLines := reBrTag.Split(normalizeString(changesFile), 3)
	for i := range splitLines {
		if reAlphaNum.MatchString(splitLines[i]) {
			a.Latest = normalizeString(splitLines[i])
			return
		}
	}
	fmt.Println(a.Name, "could not set online version, reached EOF without match")
}

func (a *addon) setLocalVersion() {
	var err error
	version, err := walkFile(makeSpecificPath(a.Path, `CHANGES.txt`), scanChanges)
	if err != nil {
		return
	}
	a.Version = normalizeString(version)
}
