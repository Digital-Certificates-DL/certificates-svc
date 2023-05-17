package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gopkg.in/gographics/imagick.v2/imagick"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

type PDF struct {
	High         float64 `json:"high"`
	Width        float64 `json:"width"`
	Name         Field   `json:"name"`
	Course       Field   `json:"course"`
	Credits      Field   `json:"credits"`
	Points       Field   `json:"points"`
	SerialNumber Field   `json:"serial_number"`
	Date         Field   `json:"date"`
	QR           Field   `json:"qr"`
	Exam         Field   `json:"exam"`
	Level        Field   `json:"level"`
	Note         Field   `json:"note"`
}

type Field struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Size  int     `json:"size"`
	Font  string  `json:"font"`
	High  float64 `json:"high"`
	Width float64 `json:"width"`
}

type PDFData struct {
	Name         string
	Course       string
	Credits      string
	Points       string
	SerialNumber string
	Date         string
	QR           []byte
	Exam         string
	Level        string
	Note         string
}

var DefaultTemplateNormal = PDF{
	High:  595,
	Width: 842,
	Name: Field{
		X:    200,
		Y:    217,
		Size: 28,
		Font: "semibold",
	},
	Course: Field{
		X:    61,
		Y:    259,
		Size: 14,
		Font: "semibold",
	},
	Credits: Field{ //todo get from front and save to db
		X:    70,
		Y:    56,
		Size: 12,
		Font: "regular",
	},
	Points: Field{
		X:    70,
		Y:    79,
		Size: 12,
		Font: "regular",
	},
	SerialNumber: Field{
		X:    641,
		Y:    56,
		Size: 12,
		Font: "regular",
	},
	Date: Field{
		X:    641,
		Y:    79,
		Size: 12,
		Font: "regular",
	},
	QR: Field{
		X:     658,
		Y:     106,
		High:  114,
		Width: 114,
	},
	Exam: Field{
		X:    300,
		Y:    300,
		Size: 15,
		Font: "italic",
	},
	Level: Field{
		X:    300,
		Y:    277,
		Size: 14,
		Font: "semibold",
	},
	//Note: Field{
	//	X:    12,
	//	Y:    466,
	//	Size: 12,
	//	Font: "arial",
	//},
}
var DefaultTemplateTall = PDF{
	High:  1190,
	Width: 1684,
	Name: Field{
		Y:    434,
		Size: 56,
		Font: "semibold",
	},
	Course: Field{
		Y:    518,
		Size: 28,
		Font: "semibold",
	},
	Credits: Field{ //todo get from front and save to db
		X:    140,
		Y:    112,
		Size: 24,
		Font: "regular",
	},
	Points: Field{
		X:    140,
		Y:    158,
		Size: 24,
		Font: "regular",
	},
	SerialNumber: Field{
		X:    1282,
		Y:    112,
		Size: 24,
		Font: "regular",
	},
	Date: Field{
		X:    1282,
		Y:    158,
		Size: 24,
		Font: "regular",
	},
	QR: Field{
		X:     1316,
		Y:     212,
		High:  228,
		Width: 228,
	},
	Exam: Field{
		Y:    600,
		Size: 30,
		Font: "italic",
	},
	Level: Field{
		Y:    554,
		Size: 28,
		Font: "semibold",
	},
}

var DefaultData = PDFData{
	Name:         "Test Name",
	Course:       "Blockchain and Distributed Systems",
	Credits:      " 99",
	Points:       "100",
	SerialNumber: "694d0f5a7afe6fbc99cb",
	Date:         "30.05.2018",
	QR:           nil,
	Exam:         "passed",
	Level:        "graduated with honors",
	Note:         "************************************************",
}

func NewPDF(high, width float64) *PDF {
	return &PDF{
		High:  high,
		Width: width,
	}
}

func NewData(name, course, credits, points, serialNumber, date string, qr []byte, exam, level, note string) PDFData {
	return PDFData{
		Name:         name,
		Course:       course,
		Credits:      credits,
		Points:       points,
		SerialNumber: serialNumber,
		Date:         date,
		QR:           qr,
		Exam:         exam,
		Level:        level,
		Note:         note,
	}
}

func (p *PDF) SetName(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Name = fl
}

func (p *PDF) SetCourse(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Course = fl
}
func (p *PDF) SetCredits(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Credits = fl
}

func (p *PDF) SetLevel(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Level = fl
}

func (p *PDF) SetPoints(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Points = fl
}

func (p *PDF) SetSerialNumber(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.SerialNumber = fl
}

func (p *PDF) SetDate(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Date = fl
}
func (p *PDF) SetQR(x, y float64, size int, high, width float64) {
	fl := Field{
		X:     x,
		Y:     y,
		Size:  size,
		High:  high,
		Width: width,
	}

	p.QR = fl
}

func (p *PDF) SetExam(x, y float64, size int, font string) {
	fl := Field{
		X:    x,
		Y:    y,
		Size: size,
		Font: font,
	}

	p.Exam = fl
}

func (p *PDF) Prepare(data PDFData, cfg config.Config, templateQ data.TemplateQ, backgroundImg []byte) ([]byte, string, []byte, error) {
	var err error
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: p.Width, H: p.High}})
	pdf.AddPage()
	pdf.SetTextColor(255, 255, 255)
	err = pdf.AddTTFFont("italic", "staff/font/Inter-Italic.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add font")
	}
	err = pdf.AddTTFFont("regular", "staff/font/Inter-Regular.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-Regular")
	}
	err = pdf.AddTTFFont("semibold", "staff/font/Inter-SemiBold.ttf")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	}
	templateImg := cfg.TemplatesConfig()[data.Course]
	//tpl1 := pdf.ImportPage(fmt.Sprintf("staff/templates/%s.pdf", templateImg), 1, "/MediaBox") //todo use  bytes

	// Draw pdf onto page
	//pdf.UseImportedTemplate(tpl1, 0, 0, 0, 0)

	if backgroundImg == nil {
		template, err := templateQ.GetByName(templateImg)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to get background img")
		}
		if template == nil {
			template, err = templateQ.GetByName("default")
			if err != nil {
				return nil, "", nil, errors.Wrap(err, "failed to get default background img")
			}
		}
		if template == nil {
			return nil, "", nil, errors.Wrap(err, "default template isn't found")
		}

		back, err := base64toJpg(template.ImgBytes)
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to decode base64 to jpeg")
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
	err = pdf.SetFont("regular", "", p.Name.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	//pdf.SetX(p.centralizeName(data.Name, p.Width, p.Name.Size))
	pdf.SetY(p.Name.Y)
	//pdf.Cell(nil, data.Name)
	pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, data.Name, gopdf.CellOption{Align: gopdf.Center})

	///////////// credits
	err = pdf.SetFont("italic", "", p.Credits.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Credits.X)
	pdf.SetY(p.Credits.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf(data.Credits))

	///////////// Points
	err = pdf.SetFont("italic", "", p.Points.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")

	}
	pdf.SetX(p.Points.X)
	pdf.SetY(p.Points.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Count of points: %s", data.Points))

	///////////// SerialNumber
	err = pdf.SetFont("italic", "", p.SerialNumber.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.SerialNumber.X)
	pdf.SetY(p.SerialNumber.Y)
	pdf.Cell(nil, data.SerialNumber)

	///////////// Date
	err = pdf.SetFont("italic", "", p.Date.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.Date.X)
	pdf.SetY(p.Date.Y)
	pdf.Cell(&gopdf.Rect{W: p.Width, H: p.High}, fmt.Sprintf("Issued on: %s", data.Date))

	///////////// Course
	err = pdf.SetFont("italic", "", p.Course.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	//pdf.SetTextColor()
	pdf.SetX(0)
	pdf.SetY(p.Course.Y)
	titles := cfg.TitlesConfig()
	isLevel, title, level := p.checkLevel(titles[templateImg])
	pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, title, gopdf.CellOption{Align: gopdf.Center})

	/////////// QR
	if data.QR != nil {
		img, _, err := image.Decode(bytes.NewReader(data.QR))
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to convert bytes to image")
		}
		err = pdf.ImageFrom(img, p.QR.X, p.QR.Y, &gopdf.Rect{W: 114, H: 114})
		if err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set image")
		}
	}

	/////////////// Exam
	err = pdf.SetFont("italic", "", p.Exam.Size)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(0)
	pdf.SetY(p.Exam.Y)
	ex := cfg.ExamsConfig()
	pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.High}, ex[data.Exam], gopdf.CellOption{Align: gopdf.Center})
	///////////// Level
	if isLevel {
		err = pdf.SetFont("italic", "", p.Level.Size)
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
		return nil, "", nil, errors.Wrap(err, "failed to  convert pdf to png")
	}
	file, err := os.Create("test.png")
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to create png file")
	}

	file.Write(imgBlob)
	file.Close()
	return pdfBlob, name, imgBlob, nil

}

func (p *PDF) ParsePoints(point string) (string, string) {
	splitedStr := strings.Split(point, "/")
	return splitedStr[0], splitedStr[1]
}

func (p *PDF) prepareLevel(level, title string) float64 {
	titleW := len(title) * p.Course.Size
	return float64(titleW - len(level)*p.Level.Size/2)
}

func (p *PDF) centralizeName(str string, width float64, size int) float64 {
	return (width/2 - (float64(size*len(str))*0.54)/2)
}

func (p *PDF) centralizeTitle(str string, width float64, size int) float64 {

	return (width/2 - (float64(size*len(str))*0.5)/2)

}
func (p *PDF) checkLevel(title string) (bool, string, string) {
	strs := strings.Split(title, "Level:")
	if len(strs) > 1 {
		return true, strs[0], fmt.Sprint("Level:", strs[1])
	}
	return false, strs[0], ""
}

func Convert(imgType string, blob []byte) ([]byte, error) {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err := mw.ReadImageBlob(blob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to  read  img blob")
	}
	mw.SetIteratorIndex(0)
	err = mw.SetImageFormat(imgType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set  format")
	}
	return mw.GetImageBlob(), nil
}

func base64toJpg(data string) ([]byte, error) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
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
