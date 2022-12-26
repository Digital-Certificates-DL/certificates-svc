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
	resultAddress := []byte("1BooKnbm48Eabw3FdPgTSudt9u4YTWKBvf")
	key := "key"


	for id, user := range users {
		msg := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		signature, _, address := Sign([]byte(key), []byte(msg))
		wantSign := resultSign
		wantAddress := resultAddress

		if  !bytes.Equal(signature, wantSign) {
			t.Errorf("got %q, wanted %q", signature, wantSign)
			continue
		}
		if  !bytes.Equal(address, wantAddress) {
			t.Errorf("got %q, wanted %q", address, wantAddress)
			continue
		}


		log.Println("PASS: ", id )

	}

}
