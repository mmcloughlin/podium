package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NameFromURL extracts the presentation name from its URL.
func NameFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	b := path.Base(u.Path)
	n := strings.TrimSuffix(b, path.Ext(b))
	return n, nil
}

var (
	images = flag.String("images", "", "directory to write slide images to")
	output = flag.String("output", "", "output filename")
)

func main() {
	flag.Parse()
	slidesURL := flag.Arg(0)

	if slidesURL == "" {
		log.Fatal("must provide presentation url")
	}

	imgs, err := Images(slidesURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("extracted %d slide images", len(imgs))

	dir := *images
	if dir == "" {
		dir, err = ioutil.TempDir("", "podium")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
	}

	format := filepath.Join(dir, "slide%03d.png")
	imageFilenames, err := WriteSlideImagesPNG(imgs, format)
	if err != nil {
		log.Fatal(err)
	}

	outputFilename := *output
	if outputFilename == "" {
		name, err := NameFromURL(slidesURL)
		if err != nil {
			log.Fatal(err)
		}
		outputFilename = name + ".pdf"
	}

	err = WriteImagesPDF(imageFilenames, outputFilename)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("written output to %s", outputFilename)
}
