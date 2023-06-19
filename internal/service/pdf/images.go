package pdf

import (
	"bytes"
	"encoding/base64"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gopkg.in/gographics/imagick.v3/imagick"
	"log"

	"image"
	"image/jpeg"
)

func Convert(imgType string, blob []byte) ([]byte, error) {
	log.Println("init")

	imagick.Initialize()
	log.Println("Terminate")

	defer imagick.Terminate()
	log.Println("NewMagickWand")

	mw := imagick.NewMagickWand()

	defer mw.Destroy()
	log.Println("ReadImageBlob")

	err := mw.ReadImageBlob(blob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to  read  img blob")
	}
	log.Println("SetIteratorIndex")

	mw.SetIteratorIndex(0)
	log.Println("SetImageFormat")

	err = mw.SetImageFormat(imgType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set  format")
	}
	log.Println("GetImageBlob")

	return mw.GetImageBlob(), nil
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
