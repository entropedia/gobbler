package main

import (
	"errors"
	"fmt"
	"strings"
	_ "time"
)

type VideoScanner struct {
}

type Video struct {
	Id     string
	Width  int
	Height int
	Depth  int
}

func (scanner VideoScanner) Scan(resource *Resource) error {
	if !strings.HasPrefix(resource.MimeType, "video/") {
		return errors.New("Unhandled mime-type")
	}

	fmt.Println("\t# Video Parser:", resource.Url)

	return nil
}
