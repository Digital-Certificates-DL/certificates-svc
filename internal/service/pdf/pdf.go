package pdf

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"image"
	"io"
	"log"
	"os"
	"strings"
)

func (p *PDF) Prepare(data PDFData, cfg config.Config, templateQ data.TemplateQ, backgroundImg []byte, userID int64) ([]byte, string, []byte, error) {
	var err error
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: p.Width, H: p.High}})
	pdf.AddPage()
	pdf.SetTextColor(255, 255, 255)
	//err = pdf.AddTTFFont("italic", "./staff/font/Inter-Italic.ttf")
	//if err != nil {
	//	return nil, "", nil, errors.Wrap(err, "failed to add font")
	//}
	//err = pdf.AddTTFFont("regular", "./staff/font/Inter-Regular.ttf")
	//if err != nil {
	//	return nil, "", nil, errors.Wrap(err, "failed to add Inter-Regular")
	//}
	//err = pdf.AddTTFFont("semibold", "./staff/font/Inter-SemiBold.ttf")
	//if err != nil {
	//	return nil, "", nil, errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	//}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(path)

	err = pdf.AddTTFFont("italic", "/usr/local/bin/staff/font/arial.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add font")
	}
	err = pdf.AddTTFFont("regular", "/usr/local/bin/staff/font/arial.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-Regular")
	}
	err = pdf.AddTTFFont("semibold", "/usr/local/bin/staff/font/arial.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	}
	templateImg := cfg.TemplatesConfig()[data.Course]

	if backgroundImg == nil {
		template, err := templateQ.GetByName(templateImg, userID)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to get background img")
		}
		if template == nil {
			titles := cfg.TitlesConfig()
			_ = titles
			template, err = templateQ.GetByName(templateImg, userID)
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "failed to get default background img")
			}

		}

		var back []byte
		if template != nil {
			back, err = base64toJpg(template.ImgBytes)
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "cant to decode img")

			}
		} else {

			file, err := os.Open(fmt.Sprintf("/usr/local/bin/staff/templates/%s.png", templateImg))
			fmt.Println(fmt.Sprintf("/usr/local/bin/stafftemplates/%s.png", templateImg))
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
	err = pdf.SetFont("regular", "", p.Name.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	//pdf.SetX(p.centralizeName(data.Name, p.Width, p.Name.FontSize))
	pdf.SetY(p.Name.Y)
	//pdf.Cell(nil, data.Name)
	fmt.Println(p.Width, p.High)
	pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, data.Name, gopdf.CellOption{Align: gopdf.Center})

	///////////// credits
	err = pdf.SetFont("italic", "", p.Credits.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Credits.X)
	pdf.SetY(p.Credits.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf(data.Credits))

	///////////// Points
	err = pdf.SetFont("italic", "", p.Points.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")

	}
	pdf.SetX(p.Points.X)
	pdf.SetY(p.Points.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Count of points: %s", data.Points))

	///////////// SerialNumber
	err = pdf.SetFont("italic", "", p.SerialNumber.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.SerialNumber.X)
	pdf.SetY(p.SerialNumber.Y)
	pdf.Cell(nil, data.SerialNumber)

	///////////// Date
	err = pdf.SetFont("italic", "", p.Date.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.Date.X)
	pdf.SetY(p.Date.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Issued on: %s", data.Date))

	///////////// Course
	err = pdf.SetFont("italic", "", p.Course.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	//pdf.SetTextColor()
	pdf.SetX(0)
	pdf.SetY(p.Course.Y)
	titles := cfg.TitlesConfig()
	isLevel, title, level := p.checkLevel(titles[templateImg])
	pdf.CellWithOption(&gopdf.Rect{W: p.Width - 50, H: p.High}, title, gopdf.CellOption{Align: gopdf.Center})

	/////////// QR
	if data.QR != nil {
		img, _, err := image.Decode(bytes.NewReader(data.QR))
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to convert bytes to image")
		}
		err = pdf.ImageFrom(img, p.QR.X, p.QR.Y, &gopdf.Rect{W: 228, H: 228})
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set image")
		}
	}

	/////////////// Exam
	err = pdf.SetFont("italic", "", p.Exam.FontSize)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(0)
	pdf.SetY(p.Exam.Y)
	ex := cfg.ExamsConfig()
	pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, ex[data.Exam], gopdf.CellOption{Align: gopdf.Center})
	///////////// Level
	if isLevel {
		err = pdf.SetFont("italic", "", p.Level.FontSize)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set font")
		}
		pdf.SetX(0)
		pdf.SetY(p.Level.Y)
		pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, level, gopdf.CellOption{Align: gopdf.Center})

	}

	parsedName := strings.Split(data.Name, " ")
	name := ""
	if len(parsedName) < 2 {
		name = fmt.Sprintf("certificate_%s_%s.pdf", parsedName[0], cfg.TemplatesConfig()[data.Course])
	} else {
		name = fmt.Sprintf("certificate_%s_%s_%s.pdf", parsedName[0], parsedName[1], cfg.TemplatesConfig()[data.Course])
	}

	pdfBlob := pdf.GetBytesPdf()

	imgBlob, err := Convert("png", pdfBlob)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to convert pdf to png")
	}
	file, err := os.Create("test.png")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to create png file")
	}

	file.Write(imgBlob)
	file.Close()
	return pdfBlob, name, imgBlob, nil

}

func (p *PDF) checkLevel(title string) (bool, string, string) {
	strs := strings.Split(title, "Level:")
	if len(strs) > 1 {
		return true, strs[0], fmt.Sprint("Level:", strs[1])
	}
	return false, strs[0], ""
}
