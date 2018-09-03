package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	images = flag.String("images", "", "directory to write slide images to")
	output = flag.String("pdf", "slides.pdf", "output pdf filename")
)

func main() {
	flag.Parse()
	url := flag.Arg(0)

	if url == "" {
		log.Fatal("must provide presentation url")
	}

	imgs, err := Images(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("extracted %d slide images", len(imgs))

	dir := *images
	if dir == "" {
		dir, err = ioutil.TempDir("", "example")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(dir)
	}

	format := filepath.Join(dir, "slide%03d.png")
	filenames, err := WriteSlideImagesPNG(imgs, format)
	if err != nil {
		log.Fatal(err)
	}

	err = WriteImagesPDF(filenames, *output)
	if err != nil {
		log.Fatal(err)
	}
}
