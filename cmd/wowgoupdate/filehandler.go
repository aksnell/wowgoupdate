package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type walkFunc func(string) (string, bool)

type FileRequest struct {
	path string
	file *os.File
}

//WalkFile takes a walkFunc and iterates over each line in the stored file until EOF or match is returned.
func (fr *FileRequest) WalkFile(wf walkFunc) string {
	if fr.file == nil {
		return ""
	}
	scanner := bufio.NewScanner(fr.file)
	for scanner.Scan() {
		if line, ismatch := wf(scanner.Text()); ismatch {
			return line
		}
	}
	return ""
}

//Close disconnects the stored os.FIle and sets it to nil.
func (fr *FileRequest) Close() {
	fr.file.Close()
	fr.file = nil
}

//Get returns a fileRequest
func Get(paths ...string) *FileRequest {
	if len(paths) < 0 {
		return nil
	}
	if !reIsExt.MatchString(paths[len(paths)-1]) {
		return nil
	}
	formattedPath := filepath.Clean(strings.Join(paths[:len(paths)-1], `\`) + paths[len(paths)-1])
	requestFile, err := os.Create(formattedPath)
	if err != nil {
		return nil
	}
	req := &FileRequest{
		path: formattedPath,
		file: requestFile,
	}
	return req
}
