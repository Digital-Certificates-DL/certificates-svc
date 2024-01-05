package pdf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/signintech/gopdf"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"image"
	"io"
	"os"
	"strconv"
	"strings"
)

func (p *PDF) Prepare(data PDFData, config *PDFConfig, masterQ data.MasterQ, backgroundImg []byte, userID int64, abs string) ([]byte, string, []byte, error) {
	pdf := new(gopdf.GoPdf)
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: p.Width, H: p.Height}})
	pdf.AddPage()
	pdf.SetTextColor(255, 255, 255)

	if err := p.setFonts(pdf, abs); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set fonts")
	}

	templateImg := config.templates[data.Course]

	if templateImg == "" {
		templateImg = data.Course
	}
	if backgroundImg == nil {
		if err := p.initBackground(pdf, masterQ.TemplateQ(), templateImg, abs, userID); err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to init background")
		}
	} else {
		if err := p.setBackground(pdf, backgroundImg); err != nil {
			return nil, "", nil, errors.Wrap(err, "failed to set background")
		}
	}

	if err := p.CellAllPdfFields(pdf, data, config, templateImg); err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to set all pdf's fields")
	}

	pdfBlob := pdf.GetBytesPdf()
	imgBlob, err := NewImageConverter().Convert(pdfBlob)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "failed to convert pdf to png")
	}

	return pdfBlob, p.prepareName(data.Name, config.templates[data.Course]), imgBlob, nil

}

func (p *PDF) checkLevel(title string) (bool, string, string) {
	titles := strings.Split(title, "Level:")
	if len(titles) > 1 {
		return true, titles[0], fmt.Sprint("Level:", titles[1])
	}
	return false, titles[0], ""
}

func (p *PDF) setBackgroundFromFile(pdf *gopdf.GoPdf, abs, imageName string) error {
	file, err := os.Open(fmt.Sprintf("%s/static/templates/%s.png", abs, imageName))
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "default template isn't found")
	}

	back, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "cant to decode img")

	}
	if err := p.setBackground(pdf, back); err != nil {
		return errors.Wrap(err, "cant to set img")
	}

	return nil
}

func (p *PDF) setBackgroundFromTemplate(pdf *gopdf.GoPdf, image []byte) error {
	back, err := NewImageConverter().base64toJpg(image)
	if err != nil {
		return errors.Wrap(err, "cant to decode img")

	}
	if err := p.setBackground(pdf, back); err != nil {
		return errors.Wrap(err, "cant to set img")
	}

	return nil
}

func (p *PDF) setBackground(pdf *gopdf.GoPdf, image []byte) error {
	backgroundImgHolder, err := gopdf.ImageHolderByBytes(image)
	if err != nil {
		return errors.Wrap(err, "failed to prepare background")
	}

	err = pdf.ImageByHolder(backgroundImgHolder, 0, 0, &gopdf.Rect{W: p.Width, H: p.Height})
	if err != nil {
		return errors.Wrap(err, "failed to set background")
	}
	return nil
}

func (p *PDF) setFonts(pdf *gopdf.GoPdf, abs string) error {
	if err := pdf.AddTTFFont("italic", abs+"/static/font/Inter-Italic.ttf"); err != nil {
		return errors.Wrap(err, "failed to add font")
	}
	if err := pdf.AddTTFFont("regular", abs+"/static/font/Inter-Regular.ttf"); err != nil {
		return errors.Wrap(err, "failed to add Inter-Regular")
	}
	if err := pdf.AddTTFFont("semibold", abs+"/static/font/Inter-SemiBold.ttf"); err != nil {
		return errors.Wrap(err, "failed to add Inter-SemiBold.ttf")
	}

	return nil
}

func (p *PDF) preparedFontStyle(isItalic bool, isBold bool) string {
	result := ""

	if isItalic {
		result += "I"
	}

	if isBold {
		result += "B"
	}

	return result
}

func (p *PDF) setLevel(pdf *gopdf.GoPdf, level string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.Level.IsItalicText, p.Level.IsBoldText), p.Level.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font Level")
	}
	pdf.SetY(p.Level.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Level.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Level.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.Level.XCenter {
		pdf.SetX(0)

		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.Height}, level, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell Level")
		}

		return nil
	}

	pdf.SetX(p.Level.X)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, level); err != nil {
		return errors.Wrap(err, "failed to cell Level")
	}

	return nil
}

func (p *PDF) setExam(pdf *gopdf.GoPdf, exam string) error {
	if err := pdf.SetFont("regular", p.preparedFontStyle(p.Exam.IsItalicText, p.Exam.IsBoldText), p.Exam.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font Exam")
	}
	pdf.SetY(p.Exam.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Exam.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Exam.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.Exam.XCenter {
		pdf.SetX(0)
		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.Height}, exam, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell Exam")
		}

		return nil
	}

	pdf.SetX(p.Exam.X)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, exam); err != nil {
		return errors.Wrap(err, "failed to cell Exam")
	}

	return nil
}

func (p *PDF) setQR(pdf *gopdf.GoPdf, qr []byte) error {
	img, _, err := image.Decode(bytes.NewReader(qr))
	if err != nil {
		return errors.Wrap(err, "failed to convert bytes to image QR")
	}

	err = pdf.ImageFrom(img, p.QR.X, p.QR.Y, &gopdf.Rect{W: p.QR.Width, H: p.QR.Height})
	if err != nil {
		return errors.Wrap(err, "failed to set image QR")
	}

	return nil
}

func (p *PDF) setCourse(pdf *gopdf.GoPdf, courseTitle string, templateImg string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.Course.IsItalicText, p.Course.IsBoldText), p.Course.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font Course")
	}

	pdf.SetY(p.Course.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Course.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Course.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if courseTitle == "" {
		courseTitle = templateImg
	}

	if p.Course.XCenter {
		pdf.SetX(0)
		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.Height}, courseTitle, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell Course")
		}
		return nil
	}

	pdf.SetX(p.Course.X)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, courseTitle); err != nil {
		return errors.Wrap(err, "failed to cell Course")
	}

	return nil
}

func (p *PDF) setSerialNumber(pdf *gopdf.GoPdf, serialNumber string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.SerialNumber.IsItalicText, p.SerialNumber.IsBoldText), p.SerialNumber.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font SerialNumber")
	}

	pdf.SetY(p.SerialNumber.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.SerialNumber.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.SerialNumber.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.SerialNumber.XCenter {
		pdf.SetX(0)

		if err := pdf.CellWithOption(&gopdf.Rect{W: 300, H: 300}, serialNumber, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell SerialNumber ")
		}
		return nil
	}

	pdf.SetX(p.SerialNumber.X)

	if err := pdf.CellWithOption(&gopdf.Rect{W: 300, H: 300}, serialNumber, gopdf.CellOption{Align: gopdf.Right}); err != nil {
		return errors.Wrap(err, "failed to cell SerialNumber ")
	}

	return nil

}

func (p *PDF) setPoints(pdf *gopdf.GoPdf, points string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.Points.IsItalicText, p.Points.IsBoldText), p.Points.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font points")

	}
	pdf.SetY(p.Points.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Points.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Points.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.Points.XCenter {
		pdf.SetX(0)

		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.Height}, fmt.Sprintf("Count of points: %s", points), gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell points")
		}

		return nil
	}

	pdf.SetX(p.Points.X)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, fmt.Sprintf("Count of points: %s", points)); err != nil {
		return errors.Wrap(err, "failed to cell points")
	}

	return nil
}

func (p *PDF) setDate(pdf *gopdf.GoPdf, date string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.Date.IsItalicText, p.Date.IsBoldText), p.Date.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font Date")
	}

	pdf.SetY(p.Date.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Date.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Date.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.Date.XCenter {
		pdf.SetX(0)

		if err := pdf.CellWithOption(&gopdf.Rect{W: 300, H: 300}, fmt.Sprintf("Issued on: %s", date), gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell Date")
		}
		return nil
	}
	pdf.SetX(p.Date.X)

	if err := pdf.CellWithOption(&gopdf.Rect{W: 300, H: 300}, fmt.Sprintf("Issued on: %s", date), gopdf.CellOption{Align: gopdf.Right}); err != nil {
		return errors.Wrap(err, "failed to cell Date")
	}
	return nil
}

func (p *PDF) setCredits(pdf *gopdf.GoPdf, credits string) error {
	if err := pdf.SetFont("italic", p.preparedFontStyle(p.Credits.IsItalicText, p.Credits.IsBoldText), p.Credits.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font credits")
	}
	pdf.SetX(p.Credits.X)
	pdf.SetY(p.Credits.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Credits.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Credits.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, fmt.Sprintf(credits)); err != nil {
		return errors.Wrap(err, "failed to cell credits")
	}

	return nil
}

func (p *PDF) setName(pdf *gopdf.GoPdf, name string) error {
	if err := pdf.SetFont("regular", p.preparedFontStyle(p.Name.IsItalicText, p.Name.IsBoldText), p.Name.FontSize); err != nil {
		return errors.Wrap(err, "failed to set font name")
	}
	pdf.SetY(p.Name.Y)
	pdf.SetTextColor(255, 255, 255)

	if p.Name.Color != "" {
		rgb, err := p.hex2RGB(strings.Replace(p.Name.Color, "#", "", 1))
		if err == nil {
			pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
		}
	}

	if p.Name.XCenter {
		pdf.SetX(0)

		if err := pdf.CellWithOption(&gopdf.Rect{W: p.Width, H: p.Height}, name, gopdf.CellOption{Align: gopdf.Center}); err != nil {
			return errors.Wrap(err, "failed to cell name")
		}
		return nil
	}

	pdf.SetX(p.Name.X)
	if err := pdf.Cell(&gopdf.Rect{W: p.Width, H: p.Height}, name); err != nil {
		return errors.Wrap(err, "failed to cell name")
	}

	return nil

}

func (p *PDF) initBackground(pdf *gopdf.GoPdf, templateQ data.TemplateQ, templateImg, abs string, userID int64) error {
	template, err := templateQ.FilterByName(templateImg).FilterByUser(userID).Get()
	if err != nil {
		return errors.Wrap(err, "failed to get background img")
	}
	if template == nil {
		if err = p.setBackgroundFromFile(pdf, abs, templateImg); err != nil {
			return errors.Wrap(err, "failed to set back  from file")
		}
	} else {
		if err = p.setBackgroundFromTemplate(pdf, template.ImgBytes); err != nil {
			return errors.Wrap(err, "failed to set back from template")

		}
	}

	return nil
}

func (p *PDF) prepareName(name, course string) string {
	parsedName := strings.Split(name, " ")
	if len(parsedName) < 2 {
		return fmt.Sprintf("certificate_%s_%s.pdf", parsedName[0], course)
	}

	return fmt.Sprintf("certificate_%s_%s_%s.pdf", parsedName[0], parsedName[1], course)
}

func (p *PDF) SetTemplateData(template PDF) *PDF {
	certificate := NewPDF(template.Height, template.Width)
	certificate.SetName(template.Name.X, template.Name.Y, template.Name.FontSize, template.Name.Font, template.Name.Color, template.Name.XCenter)
	certificate.SetDate(template.Date.X, template.Date.Y, template.Date.FontSize, template.Date.Font, template.Date.Color, template.Date.XCenter)
	//certificate.SetCourse(template.Course.X, template.Course.Y, template.Course.FontSize, template.Course.Font, template.Course.Color, template.Course.Text, template.Course.XCenter)
	certificate.SetCredits(template.Credits.X, template.Credits.Y, template.Credits.FontSize, template.Credits.Font, template.Credits.Color, template.Credits.XCenter)
	certificate.SetExam(template.Exam.X, template.Exam.Y, template.Exam.FontSize, template.Exam.Font, template.Exam.Color, template.Exam.XCenter)
	certificate.SetLevel(template.Level.X, template.Level.Y, template.Level.FontSize, template.Level.Font, template.Level.Color, template.Level.Text, template.Level.XCenter)
	certificate.SetSerialNumber(template.SerialNumber.X, template.SerialNumber.Y, template.SerialNumber.FontSize, template.SerialNumber.Font, template.SerialNumber.Color, template.SerialNumber.XCenter)
	certificate.SetPoints(template.Points.X, template.Points.Y, template.Points.FontSize, template.Points.Font, template.Points.Color, template.Points.XCenter)
	certificate.SetQR(template.QR.X, template.QR.Y, template.QR.FontSize, template.QR.Height, template.QR.Width)

	return certificate
}

func (p *PDF) InitTemplate(masterQ data.MasterQ, templateName string, userID int64) (*PDF, error) {
	template, err := masterQ.TemplateQ().FilterByName(templateName).FilterByUser(userID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get template data")
	}
	if template == nil || template.Template == nil || template.IsDefaultTemplate {
		return &DefaultTemplateTall, nil
	}

	pdf, err := p.templateDecoder(template.Template)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode template")
	}

	return pdf, nil

}

func (p *PDF) templateDecoder(templateBytes []byte) (*PDF, error) {
	pdf := new(PDF)
	err := json.Unmarshal(templateBytes, pdf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode template")
	}

	return pdf, nil
}

func (p *PDF) CellAllPdfFields(pdf *gopdf.GoPdf, data PDFData, config *PDFConfig, templateImg string) error {
	if err := p.setName(pdf, data.Name); err != nil {
		return errors.Wrap(err, "failed to set name")
	}

	if err := p.setCredits(pdf, data.Credits); err != nil {
		return errors.Wrap(err, "failed to set credits")
	}

	if err := p.setPoints(pdf, data.Points); err != nil {
		return errors.Wrap(err, "failed to set points")
	}

	if err := p.setSerialNumber(pdf, data.SerialNumber); err != nil {
		return errors.Wrap(err, "failed to set serial number")
	}

	if err := p.setDate(pdf, data.Date); err != nil {
		return errors.Wrap(err, "failed to set data")
	}

	isLevel, title, level := p.checkLevel(config.titles[templateImg])
	if title == "" {
		level = p.Level.Text
		title = p.Course.Text
		isLevel = len(level) > 0
	}

	//if err := p.setCourse(pdf, title, templateImg); err != nil {
	//	return errors.Wrap(err, "failed to set course")
	//}

	if data.QR != nil {
		if err := p.setQR(pdf, data.QR); err != nil {
			return errors.Wrap(err, "failed to set qr")
		}
	}

	if err := p.setExam(pdf, config.exams[data.Exam]); err != nil {
		return errors.Wrap(err, "failed to set exam")
	}

	if isLevel {
		if err := p.setLevel(pdf, level); err != nil {
			return errors.Wrap(err, "failed to set level")
		}
	}

	return nil
}

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (p *PDF) hex2RGB(hex string) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(hex, 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}
