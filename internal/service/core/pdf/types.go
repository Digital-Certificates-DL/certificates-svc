package pdf

const (
	readyStatus        = "ready_status"
	isProcessingStatus = "processing"
)

type PDF struct {
	Height       float64 `json:"height"`
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
	Height   float64 `json:"height"`
	Width    float64 `json:"width"`
	Text     string  `json:"text"`
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
