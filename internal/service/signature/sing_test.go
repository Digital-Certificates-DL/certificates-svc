package signature

import (
	"bytes"
	"fmt"
	"helper/internal/data"
	"log"
	"testing"
)

func TestSign(t *testing.T) {

	users := []data.User{
		data.User{
			Date:        "20.12.2022",
			Participant: "Olena Sporova",
			CourseTitle: "Theory of database organization and basic SQL",
		},
		data.User{
			Date:        "10.05.2016",
			Participant: "Nikita Magda",
			CourseTitle: "Cryptocurrencies and Distributed Systems",
		},
	}

	resultAddress := []byte("1KEXKjW8R2dLUyjY4MLLi49FM1T1tyH1rx")
	keyWif := "5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ" //KwYitv9awg9xGVHJJ2SZVdE7WXYuXp9mJvSiz8XzyiNeu34WDsXo //5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ //09c7d74ee3ee97c9101e308e3dcec263556751b1e0b78e788add115741f61786
	resultPubKey := "02a423740a5ad7500cd9d86c6a57a6f8b81619670697e86de4fe8eaac634f0330b"

	for _, user := range users {
		msg := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		fmt.Println(msg)

		signauture := NewSignature(keyWif, msg)

		sign, pub, address, err := signauture.Sign()
		if err != nil {
			log.Println(err)
		}

		log.Printf("%x", sign)

		wantAddress := resultAddress
		wantPub := resultPubKey

		if !bytes.Equal(address, wantAddress) {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}

		pubString := fmt.Sprintf("%x", pub)
		if wantPub != pubString {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}

	}

}
