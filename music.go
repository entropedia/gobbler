package main

import (
	"errors"
	"fmt"
	"github.com/vbatts/go-taglib/taglib"
	"strings"
	"time"
)

type TrackScanner struct {
}

type Artist struct {
	Id    string
	Name  string
	Names []Name
}

type Album struct {
	Id          string
	Name        string
	Names       []Name
	Artist      string
	ReleaseDate time.Time
	Tracks      []string
}

type Track struct {
	Id          string
	Name        string
	Names       []Name
	Artist      string
	ReleaseDate time.Time
	Duration    int
}

func (track Track) Debug() {
	fmt.Println("\t\t+ Audio tags:", track.Name)
}

func (scanner TrackScanner) Scan(resource *Resource) error {
	if !strings.HasPrefix(resource.MimeType, "audio/") {
		return errors.New("Unhandled mime-type")
	}

	fmt.Println("\t# Music Parser:", resource.Url)

	f := taglib.Open(resource.Url)
	if f == nil {
		return errors.New("Can't open " + resource.Url)
	}
	defer f.Close()

	tags := f.GetTags()
	props := f.GetProperties()

	if props.Length > 0 && len(tags.Title) > 0 {
		track := Track{
			Name:     tags.Title,
			Duration: props.Length,
		}
		resource.MetaData = append(resource.MetaData, track)

		//fmt.Printf("\t\t%#v\n", tags)
		//fmt.Printf("\t\t%#v\n", props)
	}

	return nil
}
