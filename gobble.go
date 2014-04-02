package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	apiHost = "api.entropedia.net"

	parsers       []ResourceInterface
	parsersInChan chan *Resource
)

type ResourceInterface interface {
	Scan(resource *Resource) error
}

type DebugInterface interface {
	Debug()
}

type Name struct {
	Name     string
	SortName string
}

type Resource struct {
	Url        string
	Content    []byte
	MimeType   string
	EntityType string
	AddedOn    time.Time
	CheckedOn  time.Time
	Sha256Sum  string
	MetaData   []interface{}
}

type ResourceRequest struct {
	Sha256   string `json:"sha256"`
	DataSize int    `json:"dataSize"`
}

type ResourcePostStruct struct {
	Resource ResourceRequest `json:"resource"`
}

func visitPath(path string, f os.FileInfo, err error) error {
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	if len(mimeType) > 0 {
		res := Resource{Url: path, MimeType: mimeType}
		parsersInChan <- &res
	}

	return nil
}

func submitResource(res *Resource) error {
	rp := ResourcePostStruct{
		Resource: ResourceRequest{
			Sha256:   res.Sha256Sum,
			DataSize: len(res.Content),
		},
	}
	b, _ := json.Marshal(rp)

	_, err := http.Post("http://"+apiHost+":8999/v1/resources", "application/json", bytes.NewBufferString(string(b)))
	return err
}

func gobbler() {
	for {
		var err error
		res := <-parsersInChan
		fmt.Printf("- Gobbling: %s\n", res.Url)

		checksum := sha256.New()
		res.Content, err = ioutil.ReadFile(res.Url)
		if err != nil {
			fmt.Println("\t! Could not sha256sum:", res.Url)
			continue
		}
		checksum.Write(res.Content)
		res.Sha256Sum = hex.EncodeToString(checksum.Sum(nil))
		fmt.Printf("\t+ Mimetype: %s, SHA256: %s\n", res.MimeType, res.Sha256Sum)

		for _, p := range parsers {
			p.Scan(res)
		}

		fmt.Println("\t# Found meta-data pieces:", len(res.MetaData))
		for _, md := range res.MetaData {
			var di DebugInterface
			switch Type := md.(type) {
			case Track:
				di = DebugInterface(Type)
			case Image:
				di = DebugInterface(Type)
			}

			di.Debug()
		}

		err = submitResource(res)
		if err != nil {
			fmt.Println("\t! Could not post resource to server:", err)
		} else {
			//      fmt.Println("Got:", resp)
		}

		fmt.Println()
	}
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if len(path) == 0 {
		fmt.Println("Usage: gobbler [path]")
		return
	}

	parsers = []ResourceInterface{}
	parsers = append(parsers, ResourceInterface(TrackScanner{}))
	parsers = append(parsers, ResourceInterface(ImageScanner{}))
	parsers = append(parsers, ResourceInterface(VideoScanner{}))
	parsers = append(parsers, ResourceInterface(TextScanner{}))

	parsersInChan = make(chan *Resource)
	go gobbler()

	err := filepath.Walk(path, visitPath)
	//FIXME: we sleep to give the gobbler thread time to finish up - nasty hack
	time.Sleep(1 * time.Second)

	if err != nil {
		fmt.Printf("! Scanning failed: %v\n", err)
	} else {
		fmt.Println("- Scanning finished successfully!")
	}
}
