package signature

import (
	"crypto/sha256"
	"fmt"
	"helper/internal/data"
)

func Hashing(user *data.User) string {
	aggregatedStr := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
	sum := sha256.Sum256([]byte(aggregatedStr))
	user.DataHash = fmt.Sprintf("%x", sum)
	if user.TxHash == "-" {
		user.SerialNumber = fmt.Sprintf("%x", sum[:10])
	}
	return string(sum[:])
}
