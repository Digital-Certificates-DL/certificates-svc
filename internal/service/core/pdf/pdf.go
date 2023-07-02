package pdf

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"image"
	"io"
	"os"
	"strings"
)

func (p *PDF) Prepare(data PDFData, config PDFConfig, masterQ data.MasterQ, backgroundImg []byte, userID int64) ([]byte, string, []byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: p.Width, H: p.High}})
	pdf.AddPage()
	pdf.SetTextColor(255, 255, 255)
	if err := pdf.AddTTFFont("italic", "/usr/local/bin/staff/font/Inter-Italic.ttf"); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add font")
	}
	if err := pdf.AddTTFFont("regular", "/usr/local/bin/staff/font/Inter-Regular.ttf"); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-Regular")
	}
	if err := pdf.AddTTFFont("semibold", "/usr/local/bin/staff/font/Inter-SemiBold.ttf"); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	}

	templateImg := config.templates[data.Course]

	if backgroundImg == nil {
		template, err := masterQ.TemplateQ().GetByName(templateImg, userID)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to get background img")
		}
		if template == nil {
			template, err = masterQ.TemplateQ().GetByName(templateImg, userID)
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "failed to get default background img")
			}

		}

		var back []byte
		if template != nil {
			back, err = NewImageConverter().base64toJpg(template.ImgBytes)
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "cant to decode img")

			}
		} else {

			file, err := os.Open(fmt.Sprintf("/usr/local/bin/staff/templates/%s.png", templateImg))
			fmt.Println(fmt.Sprintf("/usr/local/bin/staff/templates/%s.png", templateImg))
			defer file.Close()
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "default template isn't found")
			}

			back, err = io.ReadAll(file)
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "cant to decode img")

			}
		}

		backgroundImgHolder, err := gopdf.ImageHolderByBytes(back)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to prepare background")
		}

		err = pdf.ImageByHolder(backgroundImgHolder, 0, 0, &gopdf.Rect{W: p.Width, H: p.High})
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set background")
		}
	} else {
		backgroundImgHolder, err := gopdf.ImageHolderByBytes(backgroundImg)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to prepare background")
		}

		err = pdf.ImageByHolder(backgroundImgHolder, 0, 0, &gopdf.Rect{W: p.Width, H: p.High})
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set background")
		}
	}

	///////// name
	if err := pdf.SetFont("regular", "", p.Name.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font name")
	}
	pdf.SetY(p.Name.Y)
	pdf.SetX(0)
	fmt.Println(p.Width, p.High)
	if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, data.Name, gopdf.CellOption{Align: gopdf.Center}); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell name")
	}

	///////////// credits
	if err := pdf.SetFont("italic", "", p.Credits.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font credits")
	}
	pdf.SetX(p.Credits.X)
	pdf.SetY(p.Credits.Y)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf(data.Credits)); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell credits")
	}

	///////////// Points
	if err := pdf.SetFont("italic", "", p.Points.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font points")

	}
	pdf.SetX(p.Points.X)
	pdf.SetY(p.Points.Y)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Count of points: %s", data.Points)); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell points")
	}

	///////////// SerialNumber
	if err := pdf.SetFont("italic", "", p.SerialNumber.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font SerialNumber")
	}

	pdf.SetX(p.SerialNumber.X)
	pdf.SetY(p.SerialNumber.Y)
	if err := pdf.Cell(nil, data.SerialNumber); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell SerialNumber ")
	}

	///////////// Date
	if err := pdf.SetFont("italic", "", p.Date.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font Date")
	}

	pdf.SetX(p.Date.X)
	pdf.SetY(p.Date.Y)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Issued on: %s", data.Date)); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell Date")
	}
	///////////// Course
	if err := pdf.SetFont("italic", "", p.Course.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font Course")
	}
	pdf.SetX(0)
	pdf.SetY(p.Course.Y)

	isLevel, title, level := p.checkLevel(config.titles[templateImg])
	if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, title, gopdf.CellOption{Align: gopdf.Center}); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell Course")

	}
	/////////// QR
	if data.QR != nil {
		img, _, err := image.Decode(bytes.NewReader(data.QR))
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to convert bytes to image QR")
		}

		err = pdf.ImageFrom(img, p.QR.X, p.QR.Y, &gopdf.Rect{W: p.QR.High, H: p.QR.High})
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set image QR")
		}
	}

	/////////////// Exam
	if err := pdf.SetFont("italic", "", p.Exam.FontSize); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font Exam")
	}
	pdf.SetX(0)
	pdf.SetY(p.Exam.Y)

	if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, config.exams[data.Exam], gopdf.CellOption{Align: gopdf.Center}); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to cell Exam")
	}
	///////////// Level
	if isLevel {
		if err := pdf.SetFont("italic", "", p.Level.FontSize); err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set font Level")
		}
		pdf.SetX(0)
		pdf.SetY(p.Level.Y)
		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, level, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to cell Level")
		}

	}

	parsedName := strings.Split(data.Name, " ")
	if len(parsedName) < 2 {
		name := fmt.Sprintf("certificate_%s_%s.pdf", parsedName[0], config.templates[data.Course])

		pdfBlob := pdf.GetBytesPdf()

		imgBlob, err := NewImageConverter().Convert(pdfBlob)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to convert pdf to png")
		}

		return pdfBlob, name, imgBlob, nil
	}

	name := fmt.Sprintf("certificate_%s_%s_%s.pdf", parsedName[0], parsedName[1], config.templates[data.Course])
	pdfBlob := pdf.GetBytesPdf()

	imgBlob, err := NewImageConverter().Convert(pdfBlob)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to convert pdf to png")
	}

	return pdfBlob, name, imgBlob, nil

}

func (p *PDF) checkLevel(title string) (bool, string, string) {
	titles := strings.Split(title, "Level:")
	if len(titles) > 1 {
		return true, titles[0], fmt.Sprint("Level:", titles[1])
	}
	return false, titles[0], ""
}
