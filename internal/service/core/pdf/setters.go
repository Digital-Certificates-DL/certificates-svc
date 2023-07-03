package pdf

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