package main

import (
	"bytes"
	"flag"
	"fmt"
	goqr "github.com/liyue201/goqr"
	qrcode "github.com/skip2/go-qrcode"
	"image"
	"io/ioutil"
	"log"
)

func main() {
	inputObject := flag.String("input", "", "path to object to be qraeted")
	flag.Parse()

	content, err := ioutil.ReadFile(*inputObject)
	if err != nil {
		log.Fatal(err)
	}

	result := string(content)

	//	png, err := qrcode.Encode("https://example.org", qrcode.Medium, 255)
	fmt.Printf(string(len(result)))
	err = qrcode.WriteFile(result, qrcode.Low, 256, "qr.png")
	_ = err

	imgdata, err := ioutil.ReadFile("qr.png")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		fmt.Printf("image.Decode error: %v\n", err)
		return
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		fmt.Printf("Recognize failed: %v\n", err)
		return
	}
	for _, qrCode := range qrCodes {
		fmt.Printf("qrCode text: %s\n", qrCode.Payload)
	}
}

func makeQrcode(data []byte, tempName string) error {

	result := string(data)

	err := qrcode.WriteFile(result, qrcode.Low, 256, tempName)
	return err

}