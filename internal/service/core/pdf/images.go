package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ImageConvertHandler interface {
	Convert(blob []byte) ([]byte, error)
	base64toJpg(data []byte) ([]byte, error)
}

type ImageConvert struct{}

func NewImageConverter() ImageConvertHandler {
	return ImageConvert{}
}

func (i ImageConvert) Convert(blob []byte) ([]byte, error) {
	unixTime := time.Now().Unix()

	fileInput, err := os.Create(fmt.Sprintf("input%d.pdf", unixTime))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file")
	}

	defer os.Remove(fileInput.Name())

	if _, err = fileInput.Write(blob); err != nil {
		return nil, errors.Wrap(err, "failed to write data")
	}

	fileInputPath, err := filepath.Abs(fileInput.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get abs input file")
	}

	fileOutput, err := os.Create(fmt.Sprintf("output%d.png", unixTime))
	if err != nil {
		return nil, errors.Wrap(err, "failed ot create file")
	}

	defer os.Remove(fileOutput.Name())

	err = fileOutput.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close file")
	}

	fileOutputPath, err := filepath.Abs(fileOutput.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get abs output file")
	}

	cmd := exec.Command("sh", "-c", "gs -sDEVICE=png16m -dNOPAUSE -dBATCH -dSAFER -sOutputFile="+fileOutputPath+" "+fileInputPath)
	_, err = cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec convert script")
	}

	fileBlob, err := os.ReadFile(fileOutputPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file ")
	}

	return fileBlob, nil
}

func (i ImageConvert) base64toJpg(data []byte) ([]byte, error) {

	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode file")
	}

	//Encode from image format to writer
	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, errors.Wrap(err, "failed to  encode data")

	}
	return buf.Bytes(), nil
}
