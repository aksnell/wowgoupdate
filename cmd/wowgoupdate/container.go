package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type addonContainer struct {
	AddonDir  string            `json:"path"`
	Installed map[string]*addon `json:"installed"`
	Ignored   map[string]bool   `json:"ignored"`
}

func (con *addonContainer) getInstalledAddons() {
	addonFolders, err := ioutil.ReadDir(con.AddonDir)
	if err != nil {
		log.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(addonFolders))
	for _, folder := range addonFolders {
		go func(folder string) {
			defer wg.Done()
			addon := buildAddon(makeSpecificPath(con.AddonDir, folder), con.AddonDir)
			if addon != nil {
				con.Installed[addon.Name] = addon
			}
		}(folder.Name())
	}
	wg.Wait()
	fmt.Println("DONE")
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
