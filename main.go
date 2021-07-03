package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	imageWidth    = 850
	imageHeight   = 130
	boxSize       = 11
	borderColor   = color.RGBA{27, 31, 35, 255}
	borderSize    = 1
	inactiveColor = color.RGBA{235, 237, 240, 255}
	fillColors    = []color.RGBA{
		{155, 233, 168, 255},
		{64, 196, 99, 255},
		{48, 161, 78, 255},
		{33, 110, 57, 255},
	}
	months = []string{"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// Most users are not active, so this code simulate the user's inactivity
	// the lower `lazynessLevel`, the more activity you see
	lazynessLevel := rand.Intn(100) + 1
	for i := 1; i <= lazynessLevel; i++ {
		fillColors = append(fillColors, inactiveColor)
	}
}

func main() {
	filePath := "output.png"
	if len(os.Args) >= 2 {
		filePath = os.Args[1]
	}

	img := createImage()
	drawBoxes(img)
	drawMonths(&img)
	drawDays(&img)

	saveToFile(img, filePath)
}

func createImage() draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	// fill the background with white
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	return img
}

func drawBox(img draw.Image, pos image.Point) {
	borderBox := image.Rectangle{pos, image.Point{pos.X + boxSize, pos.Y + boxSize}}
	box := image.Rectangle{image.Point{pos.X + borderSize, pos.Y + borderSize},
		image.Point{pos.X + boxSize - borderSize, pos.Y + boxSize - borderSize}}

	// draw border
	draw.Draw(img, borderBox, &image.Uniform{borderColor}, image.Point{}, draw.Src)
	// fill the box
	draw.Draw(img, box, &image.Uniform{fillColors[rand.Intn(len(fillColors))]}, image.Point{}, draw.Src)
}

func drawBoxes(img draw.Image) draw.Image {
	padding := 4
	marginLeft := 50
	y := 22
	rows := 7
	cols := 53

	for j := 1; j <= rows; j++ {
		for i := 0; i < cols; i++ {
			x := (padding + boxSize) * i
			drawBox(img, image.Point{x + marginLeft, y})
		}
		y = y + boxSize + padding
	}

	return img
}

func drawString(img *draw.Image, x, y int, text string) {
	black := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := font.Drawer{
		Dst:  *img,
		Src:  image.NewUniform(black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

func drawMonths(img *draw.Image) {
	padding := 66
	x := 70
	for _, m := range months {
		drawString(img, x, 15, m)
		x += padding
	}
}

func drawDays(img *draw.Image) {
	y := 48
	for _, m := range []string{"Mon", "Wed", "Fri"} {
		drawString(img, 15, y, m)
		y += 30
	}
}

func saveToFile(img image.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Couldn't create the file")
	}
	defer f.Close()
	png.Encode(f, img)
}
