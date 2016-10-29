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
	"image/png"
	"log"
	"os"
	"os/exec"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	var textToConvert string
	if len(os.Args) <= 1 {
		textToConvert = askForInput()
		fmt.Println(textToConvert)
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
		qrcode, err = barcode.Scale(qrcode, 100, 100)
		if err != nil {
			fmt.Println(err)
		} else {
			png.Encode(f, qrcode)
		}
	}
	cmd := exec.Command("xdg-open", fileName)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

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
