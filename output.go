package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/jung-kurt/gofpdf"
)

// WriteSlideImagesPNG writes all images in PNG format, using the format to
// generate filenames (passed fmt.Sprintf). Returns list of image filenames.
func WriteSlideImagesPNG(imgs []image.Image, format string) ([]string, error) {
	var buf bytes.Buffer
	filenames := make([]string, len(imgs))
	for i, img := range imgs {
		buf.Reset()
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}

		filename := fmt.Sprintf(format, i+1)
		if err := ioutil.WriteFile(filename, buf.Bytes(), 0660); err != nil {
			return nil, err
		}

		log.Printf("slide %d image written to %s", i+1, filename)
		filenames[i] = filename
	}
	return filenames, nil
}

// WriteImagesPDF builds a PDF from slide images.
func WriteImagesPDF(filenames []string, output string) error {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "in",
		Size:    gofpdf.SizeType{Wd: 11, Ht: 7},
	})

	opt := gofpdf.ImageOptions{
		ImageType: "png",
	}

	for _, filename := range filenames {
		pdf.AddPage()
		pdf.ImageOptions(filename, 0, 0, 11, 0, false, opt, 0, "")
	}

	return pdf.OutputFileAndClose(output)
}
