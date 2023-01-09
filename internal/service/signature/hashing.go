package signature

import (
	"crypto/sha256"
	"fmt"
)

func (s Signature) Hashing() string {
	sum := sha256.Sum256([]byte(s.msg))
	return fmt.Sprintf("%x", sum[:])
}
