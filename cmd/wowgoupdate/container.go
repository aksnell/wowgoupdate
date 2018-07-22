package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type addonContainer struct {
	AddonDir  string            `json:"path"`
	Installed map[string]*addon `json:"installed"`
	Ignored   map[string]bool   `json:"ignored"`

	errHandler *log //Runtime errors
}

func (con *addonContainer) setInstalledAddons() {
	addonFolders, err := ioutil.ReadDir(con.AddonDir)
	if err != nil {
		con.log("getInstalledAddons", err, 0)
	}
	addonChannel := make(chan *addon)
	doneChannel := make(chan interface{})
	go func(numAddons int) {
		for {
			select {
			case addon := <-addonChannel:
				numAddons--
				if addon != nil {
					con.log("setInstalledAddons", errors.New("added "+addon.Name), 2)
					con.Installed[addon.Name] = addon
				}
			default:
				if numAddons <= 0 {
					doneChannel <- nil
					return
				}
			}
		}
	}(len(addonFolders))
	for _, addonFolder := range addonFolders {
		go func(addonPath string) {
			addonChannel <- buildAddon(addonPath, con.errHandler)
		}(getRelPath(con.AddonDir, addonFolder.Name()))
	}
	<-doneChannel
}

func (con *addonContainer) loadSaved() {
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		con.log(containerLoadFileErr, err, 1)
	}
	err = json.Unmarshal(data, con)
	if err != nil {
		con.log(containerLoadFileErr, err, 0)
	}
}

func (con *addonContainer) log(base string, prefix error, level int) {
	con.errHandler.add(base, prefix, level)
}
