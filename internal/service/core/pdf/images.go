package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Convert(blob []byte) ([]byte, error) {

	time := time.Now().Unix()

	fileInput, err := os.Create(fmt.Sprintf("input%d.pdf", time))
	if err != nil {
		return nil, err
	}
	_, err = fileInput.Write(blob)
	if err != nil {
		return nil, err
	}

	fileInputPath, err := filepath.Abs(fileInput.Name())
	if err != nil {
		return nil, err
	}

	fileOutput, err := os.Create(fmt.Sprintf("output%d.png", time))
	if err != nil {
		return nil, err
	}

	err = fileOutput.Close()
	if err != nil {
		return nil, err
	}

	fileOutputPath, err := filepath.Abs(fileOutput.Name())
	if err != nil {
		return nil, err
	}

	log.Println(fileInputPath + " " + fileOutputPath)
	cmd := exec.Command("sh", "-c", "gs -sDEVICE=png16m -dNOPAUSE -dBATCH -dSAFER -sOutputFile="+fileOutputPath+" "+fileInputPath)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	log.Println(out)

	fileBlob, err := os.ReadFile(fileOutputPath)
	if err != nil {
		return nil, err
	}

	return fileBlob, nil
}

func base64toJpg(data []byte) ([]byte, error) {

	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	//Encode from image format to writer
	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
