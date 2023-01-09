package signature

import (
	"crypto/sha256"
)

func (s Signature) Hashing() string {

	sum := sha256.Sum256([]byte(s.msg))

	return string(sum[:])
}
