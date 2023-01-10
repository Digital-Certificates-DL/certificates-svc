package data

type User struct {
	Date               string
	Participant        string
	CourseTitle        string
	SerialNumber       string
	Note               string
	Certificate        string
	DataHash           string
	TxHash             string
	Signature          string
	DigitalCertificate string
	ID                 int
}

func (u *User) SetSignature(signature string) {
	u.Signature = signature
}

func (u *User) SetDataHash(hash string) {
	u.DataHash = hash
	if u.TxHash == "-" {
		u.TxHash = hash[:10]
	}

}
