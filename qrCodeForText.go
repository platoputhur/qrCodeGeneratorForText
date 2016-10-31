/*Copyright (c) 2016 Plato Puthur. All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
)

func main() {
	var textToConvert string
	if len(os.Args) <= 1 {
		textToConvert = askForInput()
	} else {
		textSlice := os.Args[1:]
		textToConvert = textSlice[0]
	}
	qrcode, err := qr.Encode(textToConvert, qr.H, qr.Unicode)

	f, _ := os.Create("qrcode.png")
	fileName := f.Name()
	defer f.Close()

	if err != nil {
		fmt.Println(err)
	} else {
		qrcode, err = barcode.Scale(qrcode, 250, 250)
		if err != nil {
			fmt.Println(err)
		} else {
			png.Encode(f, qrcode)
		}
	}
	var i imageViewerConfig
	i.createImageFileWithBorder()
	gl.StartDriver(i.createImageWindows)

	fmt.Println("Delete the generated qr code?")
	if askForConfirmation() == true {
		err = os.Remove(fileName)
		fmt.Println("Qr code deleted!")
	} else {
		fmt.Println("Qr code not deleted, can be found in the working directory!")
	}
}

func askForInput() string {
	fmt.Println("You didn't give me an input, please enter the text you would like to convert to QR Code:")
	var textToConvert string
	_, err := fmt.Scanln(&textToConvert)
	if err != nil {
		textToConvert = askForInput()
	}
	return textToConvert
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func (i *imageViewerConfig) createImageFileWithBorder() {
	imageFile, err := os.Open("qrcode.png")

	if err != nil {
		fmt.Println(err)
	}

	sourceImage, _, err := image.Decode(imageFile)
	imageFile.Close()
	i.m = image.NewRGBA(image.Rect(0, 0, sourceImage.Bounds().Dx()+6, sourceImage.Bounds().Dy()+6))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(i.m, i.m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	if err != nil {
		fmt.Println(err)
	}

	imageFile, _ = os.Create("qrcode.png")
	defer imageFile.Close()

	qrSourceImage := sourceImage.Bounds()
	draw.Draw(i.m, qrSourceImage.Add(image.Point{3, 3}), sourceImage, image.ZP, draw.Src)
	png.Encode(imageFile, i.m)
}

func (i *imageViewerConfig) createImageWindows(driver gxui.Driver) {
	width, height := 256, 256
	theme := dark.CreateTheme(driver)
	img := theme.CreateImage()
	window := theme.CreateWindow(width, height, "QrCode viewer")
	texture := driver.CreateTexture(i.m, 1.0)
	img.SetTexture(texture)
	window.AddChild(img)
	window.OnClose(driver.Terminate)
}

type imageViwer interface {
	createImageWindows()
	createImageFileWithBorder()
}

type imageViewerConfig struct {
	m           *image.RGBA
	sourceImage *image.RGBA
}
