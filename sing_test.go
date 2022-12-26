package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

package main

import (
"fmt"
"log"
"testing"
)

func TestSign(t *testing.T) {

	users := []user{
		user{
			Date:        "20.12.2022",
			Participant: "Olena Sporova",
			CourseTitle: "Theory of database organization and basic SQL",
		},

		user{
			Date:        "10.05.2016",
			Participant: "Nikita Magda",
			CourseTitle: "Cryptocurrencies and Distributed Systems",
		},
	}

	resultSign := []byte("HAQnJ3w8TAAUbuNJ7rcf8pjR1ef1LZJXme8DC6ORv5jTCW/tNlrDCYagQUv+t7oqZtRj/3KhIVV8rgZ2txsysPg=")
	key := "key"
	for _, user := range users {
		msg := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		signature, _, _ := Sign([]byte(key), []byte(msg))
		want := resultSign

		if  bytes.Equal(signature, want) {
			t.Errorf("got %q, wanted %q", signature, want)
			continue
		}
		log.Println("PASS")

	}

}
