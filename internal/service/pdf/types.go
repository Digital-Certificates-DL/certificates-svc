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

var DefaultTemplateNormal = PDF{
	High:  595,
	Width: 842,
	Name: Field{
		X:        200,
		Y:        217,
		FontSize: 28,
		Font:     "semibold",
	},
	Course: Field{
		X:        61,
		Y:        259,
		FontSize: 14,
		Font:     "semibold",
	},
	Credits: Field{
		X:        70,
		Y:        56,
		FontSize: 12,
		Font:     "regular",
	},
	Points: Field{
		X:        70,
		Y:        79,
		FontSize: 12,
		Font:     "regular",
	},
	SerialNumber: Field{
		X:        641,
		Y:        56,
		FontSize: 12,
		Font:     "regular",
	},
	Date: Field{
		X:        641,
		Y:        79,
		FontSize: 12,
		Font:     "regular",
	},
	QR: Field{
		X:     658,
		Y:     106,
		High:  114,
		Width: 114,
	},
	Exam: Field{
		X:        300,
		Y:        300,
		FontSize: 15,
		Font:     "italic",
	},
	Level: Field{
		X:        300,
		Y:        277,
		FontSize: 14,
		Font:     "semibold",
	},
}

var DefaultTemplateTall = PDF{
	High:  1190,
	Width: 1684,
	Name: Field{
		Y:        434,
		FontSize: 56,
		Font:     "semibold",
	},
	Course: Field{
		Y:        518,
		FontSize: 28,
		Font:     "semibold",
	},
	Credits: Field{ //todo get from front and save to db
		X:        140,
		Y:        112,
		FontSize: 24,
		Font:     "regular",
	},
	Points: Field{
		X:        140,
		Y:        158,
		FontSize: 24,
		Font:     "regular",
	},
	SerialNumber: Field{
		X:        1144,
		Y:        112,
		FontSize: 24,
		Font:     "regular",
	},
	Date: Field{
		X:        1282,
		Y:        158,
		FontSize: 24,
		Font:     "regular",
	},
	QR: Field{
		X:     1316,
		Y:     212,
		High:  228,
		Width: 228,
	},
	Exam: Field{
		Y:        600,
		FontSize: 30,
		Font:     "italic",
	},
	Level: Field{
		Y:        554,
		FontSize: 28,
		Font:     "semibold",
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
