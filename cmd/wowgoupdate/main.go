package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	config = "config.json"
)

func main() {
	log := &log{err: make(map[int][]string)}
	container, err := buildContainer(log)
	if err != nil {
		log.add("buildContainer could not build", err, 0)
	}
	container.setInstalledAddons()
	save(container)
	log.dump(critical)
	for {
	}
}

func buildContainer(logger *log) (*addonContainer, error) {
	addonsDir, err := walkFile(productFile, scanProduct)
	if err != nil {
		return nil, err
	}
	container := &addonContainer{
		AddonDir:   addonsDir,
		Installed:  make(map[string]*addon),
		Ignored:    make(map[string]bool),
		errHandler: logger,
	}
	return container, nil
}

func save(con *addonContainer) {
	file, _ := os.OpenFile(
		saveFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	defer file.Close()
	data, err := json.MarshalIndent(con, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
