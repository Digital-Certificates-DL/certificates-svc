package pdf

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
)

//func Convert(imgType string, blob []byte) ([]byte, error) {
//	log.Println("init")
//
//	imagick.Initialize()
//	log.Println("Terminate")
//
//	defer imagick.Terminate()
//	log.Println("NewMagickWand")
//
//	mw := imagick.NewMagickWand()
//
//	defer mw.Destroy()
//	log.Println("ReadImageBlob")
//
//	err := mw.ReadImageBlob(blob)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to  read  img blob")
//	}
//	log.Println("SetIteratorIndex")
//
//	mw.SetIteratorIndex(0)
//	log.Println("SetImageFormat")
//
//	err = mw.SetImageFormat(imgType)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to set  format")
//	}
//	log.Println("GetImageBlob")
//
//	return mw.GetImageBlob(), nil
//}

func Convert(imgType string, blob []byte) ([]byte, error) {
	fileInput, err := os.Create("input.pdf")
	if err != nil {
		return nil, err
	}
	_, err = fileInput.Write(blob)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("gs -sDEVICE=png16m -dNOPAUSE -dBATCH -dSAFER -sOutputFile=output.png input.pdf")
	_, err = cmd.Output()
	if err != nil {
		return nil, err
	}

	fileBlob, err := os.ReadFile("output.png")
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
