package grata

func Enqrata(file *File) error {
	dir := createTmpDir()

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
