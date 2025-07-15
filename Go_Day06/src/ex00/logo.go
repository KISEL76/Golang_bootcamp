package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	const width, height = 300, 300

	bgColor := color.RGBA{10, 20, 60, 255}      // navy blue
	letterColor := color.RGBA{220, 20, 60, 255} // red

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	centerX := width/2 + 25
	topY := 50
	bottomY := 250

	thickness := 10

	// verical line "K"
	for x := centerX - 65; x < centerX-60+thickness; x++ {
		for y := topY; y < bottomY; y++ {
			img.Set(x, y, letterColor)
		}
	}

	// diagonal up
	for i := 0; i < 100; i++ {
		for t := -thickness / 2; t <= thickness/2; t++ {
			x := centerX - 60 + i/2 + t
			y := 150 - i
			if x >= 0 && x < width && y >= 0 && y < height {
				img.Set(x, y, letterColor)
			}
		}
	}

	// diagonal down
	for i := 0; i < 100; i++ {
		for t := -thickness / 2; t <= thickness/2; t++ {
			x := centerX - 60 + i/2 + t
			y := 150 + i
			if x >= 0 && x < width && y >= 0 && y < height {
				img.Set(x, y, letterColor)
			}
		}
	}

	// Save
	file, err := os.Create("amazing_logo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
