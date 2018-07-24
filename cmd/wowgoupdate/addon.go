package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type addon struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Latest      string `json:"latest"`
	DownloadURL string `json:"download"`
}

func buildAddon(path string, addonDir string) *addon {
	addon := &addon{
		Path: path,
	}
	addon.setName()
	if addon.Name == "" {
		return nil
	}
	addon.setURL()
	if addon.URL == "" {
		return nil
	}
	addon.setDownloadURL()
	if addon.DownloadURL == "" {
		return nil
	}
	addon.setOnlineVersion()
	if addon.Latest == "" {
		return nil
	}
	addon.setLocalVersion()
	addon.setLocalVersionBrute()
	resp, _ := http.Get(addon.DownloadURL)
	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("unzipping:", addon.Name, "version:", addon.Latest, "to:", addon.Path)
	err := Unzip(bytes.NewReader(r), resp.ContentLength, addonDir)
	if err != nil {
		fmt.Println(err)
	}
	return addon
}

func (a *addon) setDownloadURL() {
	res, _ := http.Get(a.URL + `/files`)
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if reDownloadURL.Match(scanner.Bytes()) {
			a.DownloadURL = a.URL + `/download/` + reDownloadURL.FindStringSubmatch(scanner.Text())[1] + `/file`
			return
		}
	}
}

func (a *addon) setName() {
	var err error
	a.Name, err = walkFile(makeGenericPath(a.Path, `.toc`), scanToc)
	if err != nil {
		return
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
}

func (a *addon) setOnlineVersion() {
	res, _ := http.Get(a.URL + `/files`)
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if reVersionNum.Match(scanner.Bytes()) {
			a.Latest = reVersionNum.FindStringSubmatch(scanner.Text())[1]
			return
		}
	}
}

func (a *addon) setLocalVersion() {
	var err error
	version, err := walkFile(makeSpecificPath(a.Path, `CHANGES.txt`), scanChanges)
	if err != nil {
		return
	}
	a.Version = normalizeString(version)
}

func (a *addon) setLocalVersionBrute() {
	group := func(ctx context.Context) (string, error) {
		f, ctx := errgroup.WithContext(ctx)
		dir, _ := ioutil.ReadDir(a.Path)
		for _, folder := range dir {
			folder := folder
			f.Go(func() error {
				f, err := os.Open(makeSpecificPath(a.Path, folder.Name()))
				if err != nil {
					return nil
				}
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), a.Latest) {
						a.Version = scanner.Text()
						return errors.New("done")
					}
					return nil
				}
				return nil
			})
		}
		return "", nil
	}
	_, _ = group(context.Background())
}

func Unzip(src io.ReaderAt, len int64, dest string) error {
	r, err := zip.NewReader(src, len)
	if err != nil {
		return err
	}
	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
