package helpers

import (
	"crypto/sha256"
	"fmt"
)

type Certificate struct {
	Date               string
	Participant        string
	CourseTitle        string
	Points             string
	Note               string
	SerialNumber       string
	Certificate        string
	DataHash           string
	TxHash             string
	Signature          string
	DigitalCertificate string
	ID                 int
	Msg                string
	ImageCertificate   []byte
	ShortCourseName    string
}

func (u *Certificate) SetSignature(signature string) {
	u.Signature = signature
}

func (u *Certificate) SetDataHash(hash string) {
	u.DataHash = hash
	if len(u.TxHash) > 0 && len(u.TxHash) < 5 {
		u.SerialNumber = hash[:20]
		return
	}
	if u.TxHash != "" {
		u.SerialNumber = u.TxHash[:20]
	}
}

func (u *Certificate) Hashing(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", sum[:])
}
