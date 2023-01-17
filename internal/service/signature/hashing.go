package signature

import (
	"crypto/sha256"
	"fmt"
)

func (s Signature) Hashing(msg string) string {
	sum := sha256.Sum256([]byte(msg))
	return fmt.Sprintf("%x", sum[:])
}
