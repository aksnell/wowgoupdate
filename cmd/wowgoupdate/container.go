package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type addonContainer struct {
	AddonDir  string            `json:"path"`
	Installed map[string]*addon `json:"installed"`
	Ignored   map[string]bool   `json:"ignored"`
}

func (con *addonContainer) setInstalledAddons() {
	addonFolders, err := ioutil.ReadDir(con.AddonDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	addonChannel := make(chan *addon)
	doneChannel := make(chan interface{})
	go func(numAddons int) {
		for {
			select {
			case addon := <-addonChannel:
				numAddons--
				if addon != nil {
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
			addon, err := buildAddon(addonPath)
			if addon == nil || err != nil {
				addonChannel <- nil
				return
			}
			addonChannel <- addon
		}(getRelPath(con.AddonDir, addonFolder.Name()))
	}
	<-doneChannel
}

func (con *addonContainer) loadSaved() {
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, con)
	if err != nil {
		fmt.Println(err)
	}
}
