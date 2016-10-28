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
	textToConvert := os.Args[1:]
	qrcode, err := qr.Encode(textToConvert[0], qr.H, qr.Unicode)

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
