package pdf

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"image"
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

var DefaultTemplate = PDF{
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
		X:    572,
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
		Y:    350,
		Size: 12,
		Font: "arial",
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

//
//func (p *PDF) SetNote(x, y float64, size int, font string) {
//	fl := Field{
//		X:    x,
//		Y:    y,
//		Size: size,
//		Font: font,
//	}
//
//	p.Note = fl
//}

func (p *PDF) Prepare(data PDFData, cfg config.Config) ([]byte, string, error) {
	var err error

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: p.Width, H: p.High}})
	pdf.AddPage()

	pdf.SetTextColor(255, 255, 255)
	err = pdf.AddTTFFont("arial", "staff/font/arial.ttf")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to add font")
	}

	err = pdf.AddTTFFont("italic", "staff/font/Inter-Italic.ttf")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to add font")
	}
	err = pdf.AddTTFFont("regular", "staff/font/Inter-Regular.ttf")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to add Inter-Regular")
	}
	err = pdf.AddTTFFont("semibold", "staff/font/Inter-SemiBold.ttf")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	}
	templateImg := cfg.TemplatesConfig()[data.Course]
	tpl1 := pdf.ImportPage(fmt.Sprintf("staff/templates/%s.pdf", templateImg), 1, "/MediaBox") //todo use  bytes

	// Draw pdf onto page
	pdf.UseImportedTemplate(tpl1, 0, 0, 0, 0)

	// Color the page
	//pdf.SetLineWidth(0.1)

	///////// name
	err = pdf.SetFont(p.Name.Font, "", p.Name.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Name.X)
	pdf.SetY(p.Name.Y)
	pdf.Cell(nil, data.Name)
	fmt.Println("set")
	///////////// Course
	err = pdf.SetFont(p.Course.Font, "", p.Course.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Course.X)
	pdf.SetY(p.Course.Y)
	pdf.Cell(nil, data.Course)

	///////////// credits
	err = pdf.SetFont(p.Credits.Font, "", p.Credits.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Credits.X)
	pdf.SetY(p.Credits.Y)
	pdf.Cell(nil, fmt.Sprintf(data.Credits))

	///////////// Points
	err = pdf.SetFont(p.Points.Font, "", p.Points.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")

	}
	pdf.SetX(p.Points.X)
	pdf.SetY(p.Points.Y)
	pdf.Cell(nil, fmt.Sprintf("Count of points: %s", data.Points))

	///////////// SerialNumber
	err = pdf.SetFont(p.SerialNumber.Font, "", p.SerialNumber.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.SerialNumber.X)
	pdf.SetY(p.SerialNumber.Y)
	pdf.Cell(nil, data.SerialNumber)

	///////////// Date
	err = pdf.SetFont(p.Date.Font, "", p.Date.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}

	pdf.SetX(p.Date.X)
	pdf.SetY(p.Date.Y)
	pdf.Cell(nil, fmt.Sprintf("Issued on: %s", data.Date))

	///////////// Course
	err = pdf.SetFont(p.Course.Font, "", p.Course.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	//pdf.SetTextColor()
	pdf.SetX(p.Course.X)
	pdf.SetY(p.Course.Y)
	pdf.Cell(nil, fmt.Sprintf("has successfully completed Course \"%s\"", data.Course))

	///////////// QR
	img, _, err := image.Decode(bytes.NewReader(data.QR))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to convert bytes to image")
	}
	err = pdf.ImageFrom(img, p.QR.X, p.QR.Y, &gopdf.Rect{W: 100, H: 100})
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set image")
	}
	///////////// Exam
	err = pdf.SetFont(p.Exam.Font, "", p.Exam.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Exam.X)
	pdf.SetY(p.Exam.Y)
	pdf.Cell(nil, fmt.Sprintf("Exam passed %s", data.Exam))

	///////////// Level
	err = pdf.SetFont(p.Level.Font, "", p.Level.Size)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to set font")
	}
	pdf.SetX(p.Level.X)
	pdf.SetY(p.Level.Y)
	pdf.Cell(nil, fmt.Sprintf("Level: %s", data.Level))

	parsedName := strings.Split(data.Name, " ")
	name := ""
	if len(parsedName) < 2 {
		name = fmt.Sprintf("certificate_%s_%s.pdf", parsedName[0], cfg.TemplatesConfig()[data.Course])
	} else {
		name = fmt.Sprintf("certificate_%s_%s_%s.pdf", parsedName[0], parsedName[1], cfg.TemplatesConfig()[data.Course])
	}
	return pdf.GetBytesPdf(), name, nil
}

func (p *PDF) ParsePoints(point string) (string, string) {
	splitedStr := strings.Split(point, "/")
	return splitedStr[0], splitedStr[1]
}

//func (p *PDF) PDFToImg(pdfData []byte) ([]byte, error) {
//	imagick.Initialize()
//	defer imagick.Terminate()
//	mw := imagick.NewMagickWand()
//	defer mw.Destroy()
//	err := mw.ReadImageBlob(pdfData)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to read pdf")
//	}
//	mw.SetIteratorIndex(0) // This being the page offset
//	err = mw.SetImageFormat("jpg")
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to set image format")
//	}
//	image := mw.GetImageBlob()
//	if image != nil {
//		return image, nil
//	}
//	return image, errors.New("failed to get image blob")
//}
