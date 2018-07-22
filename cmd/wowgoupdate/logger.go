package main

import (
	"fmt"
	"strings"
	"time"
)

type log struct {
	err    map[int][]string
	target string
}

func (l *log) add(prefix string, err error, level int) {
	formattedError := strings.Join([]string{time.Now().Format("2006-01-02 3:4 PM"), prefix, err.Error()}, ":")
	l.err[level] = append(l.err[level], formattedError)
}

func (l *log) dump(dumpToLevel int) {
	for level := 0; level < 3; level++ {
		for _, err := range l.err[level] {
			fmt.Println(err)
		}
	}
}
