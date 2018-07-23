//go:binary-only-package
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	container, err := buildContainer()
	if err != nil {
		fmt.Println(err)
	}
	container.getInstalledAddons()
	err = save(container)
	if err != nil {
		fmt.Println(err)
	}
	for {

	}
}

func buildContainer() (*addonContainer, error) {
	addonsDir, err := walkFile(productFile, scanProduct)
	if err != nil {
		return nil, err
	}
	container := &addonContainer{
		AddonDir:  addonsDir,
		Installed: make(map[string]*addon),
		Ignored:   make(map[string]bool),
	}
	return container, nil
}

func save(con *addonContainer) error {
	saveData, err := json.MarshalIndent(con, "", "	")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(saveFile, saveData, 0644); err != nil {
		return err
	}
	return nil
}
