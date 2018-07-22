package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type walkFunc func(string) (string, bool)

func makeFilePath(root string, ext string) string {
	return strings.Join([]string{root, filepath.Base(root) + ext}, `\`)
}
func getRelPath(root string, base string) string {
	return filepath.Join(root, base)
}

func walkText(text string, wF walkFunc) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		if line, err := wF(scanner.Text()); err == true {
			return line, nil
		}
	}
	return "", errors.New("walkText:Reached EOT while scanning text")
}

func walkFile(path string, wF walkFunc) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.New("walkFile:Could not load file from " + path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line, err := wF(scanner.Text()); err == true {
			return line, nil
		}
	}
	return "", errors.New("walkFile:Reached EOF while scanning " + path)
}

func scanProduct(line string) (string, bool) {
	const suffix string = `\Interface\AddOns`
	if ismatch := reAddonFolderPath.FindString(line); ismatch != "" {
		path := strings.Replace(ismatch, `/`, `\`, -1)
		return path + suffix, true
	}
	return "", false
}

func scanToc(line string) (string, bool) {
	ismatch := reAddonName.FindStringSubmatch(line)
	if len(ismatch) > 1 && ismatch[1] != "" {
		return strings.Trim(strings.Join(ismatch[1:len(ismatch)], ""), " "), true
	}
	return "", false
}

func scanChanges(line string) (string, bool) {
	if ismatch := reAlphaNum.MatchString(line); ismatch {
		return strings.Trim(line, " "), true
	}
	return "", false
}

func scanChangesURL(line string) {

}
