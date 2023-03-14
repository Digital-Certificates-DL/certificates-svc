package data

type User struct {
	Date               string
	Participant        string
	CourseTitle        string
	Points             string
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
	if len(u.TxHash) > 0 && len(u.TxHash) < 5 {
		u.SerialNumber = hash[:20]
	}

}
