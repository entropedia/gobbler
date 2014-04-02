package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type TextScanner struct {
}

type Text struct {
	Id          string
	Author      string
	ReleaseDate time.Time
}

func (scanner TextScanner) Scan(resource *Resource) error {
	if !strings.HasPrefix(resource.MimeType, "text/") {
		return errors.New("Unhandled mime-type")
	}

	fmt.Println("\t# Text Parser:", resource.Url)
	return nil
}
