package main

import (
	"flag"
	"log"
	"path/filepath"
)

var (
	images = flag.String("images", "", "directory to write slide images to")
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

	if *images != "" {
		format := filepath.Join(*images, "slide%03d.png")
		if err := WriteSlideImagesPNG(imgs, format); err != nil {
			log.Fatal(err)
		}
	}
}
