package main

import (
	"crypto/sha256"
	"fmt"
)

func hashing(user *user) string {
	aggregatedStr := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
	sum := sha256.Sum256([]byte(aggregatedStr))
	user.DataHash = fmt.Sprintf("%x", sum)
	return string(sum[:])
}
