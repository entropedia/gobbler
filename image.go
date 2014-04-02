package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	_ "time"
)

type ImageScanner struct {
}

type Image struct {
	Id     string
	Width  int
	Height int
	Depth  int
}

func (img Image) Debug() {
	fmt.Printf("\t\t+ Image resolution: %dx%dpx\n", img.Width, img.Height)
}

func (scanner ImageScanner) Scan(resource *Resource) error {
	if !strings.HasPrefix(resource.MimeType, "image/") {
		return errors.New("Unhandled mime-type")
	}

	fmt.Println("\t# Image Parser:", resource.Url)

	fo, err := os.Create("/tmp/entropedia_gobble.tmp")
	if err != nil {
		return errors.New("Could not create temporary image file")
	}
	defer fo.Close()

	_, err = fo.Write(resource.Content)
	if err != nil {
		return errors.New("Could not write to temporary image file")
	}
	fo.Close()

	fi, _ := os.Open("/tmp/entropedia_gobble.tmp")
	defer fi.Close()

	img, _, err := image.Decode(fi)
	if err != nil {
		return errors.New("Could not decode image temporary image file")
	}
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y

	if width > 0 && height > 0 {
		i := Image{
			Width:  width,
			Height: height,
		}
		resource.MetaData = append(resource.MetaData, i)
	}

	return nil
}
