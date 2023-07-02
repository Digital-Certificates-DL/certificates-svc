package helpers

import (
	"crypto/sha256"
	"fmt"
)

type User struct {
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

func (u *User) SetSignature(signature string) {
	u.Signature = signature
}

func (u *User) SetDataHash(hash string) {
	u.DataHash = hash
	if len(u.TxHash) > 0 && len(u.TxHash) < 5 {
		u.SerialNumber = hash[:20]
		return
	}
	if u.TxHash != "" {
		u.SerialNumber = u.TxHash[:20]
	}
}

func (u *User) Hashing(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", sum[:])
}
