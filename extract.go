package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/oliamb/cutter"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// FindSlideBounds discovers the bounding rectangle of the slide within a
// screenshot of the browser viewport.
func FindSlideBounds(img image.Image) image.Rectangle {
	bounds := img.Bounds()
	mid := (bounds.Min.X + bounds.Max.X) / 2
	top := 1
	for img.At(mid, bounds.Min.Y+top) == img.At(mid, bounds.Min.Y) {
		top++
	}

	slideHeight := bounds.Dy() - 2*top
	slideWidth := (slideHeight * 11) / 7
	left := (bounds.Dx() - slideWidth) / 2

	return image.Rectangle{
		Min: image.Point{X: bounds.Min.X + left, Y: bounds.Min.Y + top},
		Max: image.Point{X: bounds.Max.X - left, Y: bounds.Max.Y - top},
	}
}

// CropSlide crops the slide out of a screenshot of the browser viewport.
func CropSlide(img image.Image) (image.Image, error) {
	b := FindSlideBounds(img)
	return cutter.Crop(img, cutter.Config{
		Width:  b.Dx() - 1,
		Height: b.Dy() - 1,
		Mode:   cutter.Centered,
	})
}

// config fetches a key from the environment, returning d if not found.
func config(k, d string) string {
	v := os.Getenv("PODIUM_" + k)
	if v == "" {
		return d
	}
	return v
}

// Images extracts slide images from a talk URL.
func Images(url string) ([]image.Image, error) {
	const port = 1337

	selenium.SetDebug(true)

	service, err := selenium.NewChromeDriverService(config("CHROME_DRIVER_PATH", "chromedriver"), port)
	if err != nil {
		return nil, err
	}
	defer service.Stop()

	cap := selenium.Capabilities{}
	cap.AddChrome(chrome.Capabilities{
		Path: config("CHROME_PATH", ""),
		Args: []string{
			"headless",
			"window-size=1100x700",
		},
	})

	urlPrefix := fmt.Sprintf("http://localhost:%d/wd/hub", port)
	wd, err := selenium.NewRemote(cap, urlPrefix)
	if err != nil {
		return nil, err
	}
	defer wd.Quit()

	err = wd.Get(url)
	if err != nil {
		return nil, err
	}

	script := `
	document.querySelector('body').style.backgroundImage = 'none';
	document.querySelector('body').style.backgroundColor = 'green';
	document.querySelectorAll('.slides article').forEach(function(s) {
		s.style.borderRadius = 0;
		s.style.border = 'none';
	})
	document.querySelector('#help').style.display = 'none';
	`
	_, err = wd.ExecuteScript(script, nil)
	if err != nil {
		return nil, err
	}

	elems, err := wd.FindElements(selenium.ByCSSSelector, ".slides article")
	if err != nil {
		return nil, err
	}

	n := len(elems)
	imgs := make([]image.Image, n)

	for i := 0; i < n; i++ {
		log.Printf("processing slide %d\n", i+1)

		// Screenshot
		b, err := wd.Screenshot()
		if err != nil {
			return nil, err
		}

		buf := bytes.NewBuffer(b)
		img, err := png.Decode(buf)
		if err != nil {
			return nil, err
		}

		imgs[i], err = CropSlide(img)
		if err != nil {
			return nil, err
		}

		// Move to the next slide
		err = wd.KeyDown(selenium.RightArrowKey)
		if err != nil {
			return nil, err
		}

		time.Sleep(time.Second)
	}

	return imgs, nil
}
