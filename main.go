package main

import (
	"bytes"
	"flag"
	"fmt"
	goqr "github.com/liyue201/goqr"
	qrcode "github.com/skip2/go-qrcode"
	"image"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	inputObject := flag.String("input", "", "path to object to be qraeted")
	flag.Parse()

	dir := createTmpDir()
	print(dir + "\n")

	file, err := os.Open(*inputObject)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 1024) // define your buffer size here.

	part := 0

	for {
		n, err := file.Read(buf)
		part += 1

		if n > 0 {
			partPath := fmt.Sprintf("%s/%d.png", dir, part)

			makeQrcode(buf, partPath)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read %d bytes: %v", n, err)
			break
		}

	}

	makeVideo(dir)

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

	os.RemoveAll(dir)
}

func makeQrcode(data []byte, tempName string) error {

	result := string(data)

	err := qrcode.WriteFile(result, qrcode.Low, 256, tempName)
	return err
}

func createTmpDir() string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrst"
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	dir, err := ioutil.TempDir("/tmp/", output.String())
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func makeVideo(dir string) {
	cmdName := "ffmpeg"
	qrPath := fmt.Sprintf("%s/%%d.png", dir)
	args := []string{
		"-i",
		qrPath,
		"./output.mp4"}
	fmt.Println(args)
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}

}
