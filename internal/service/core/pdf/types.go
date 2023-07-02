package pdf

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
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	XCenter  bool    `json:"x_center"`
	YCenter  bool    `json:"y_center"`
	FontSize int     `json:"font_size"`
	Color    string  `json:"color"`
	Font     string  `json:"font"`
	High     float64 `json:"high"`
	Width    float64 `json:"width"`
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
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Name = fl
}

func (p *PDF) SetCourse(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Course = fl
}
func (p *PDF) SetCredits(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Credits = fl
}

func (p *PDF) SetLevel(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Level = fl
}

func (p *PDF) SetPoints(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Points = fl
}

func (p *PDF) SetSerialNumber(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.SerialNumber = fl
}

func (p *PDF) SetDate(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Date = fl
}
func (p *PDF) SetQR(x, y float64, size int, high, width float64) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		High:     high,
		Width:    width,
	}

	p.QR = fl
}

func (p *PDF) SetExam(x, y float64, size int, font string) {
	fl := Field{
		X:        x,
		Y:        y,
		FontSize: size,
		Font:     font,
	}

	p.Exam = fl
}
