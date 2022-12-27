package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
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

	resultAddress := []byte("1KEXKjW8R2dLUyjY4MLLi49FM1T1tyH1rx")
	keyWif := "5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ" //KwYitv9awg9xGVHJJ2SZVdE7WXYuXp9mJvSiz8XzyiNeu34WDsXo //5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ //09c7d74ee3ee97c9101e308e3dcec263556751b1e0b78e788add115741f61786
	key := []byte("09c7d74ee3ee97c9101e308e3dcec263556751b1e0b78e788add115741f61786")
	resultPubKey := "02a423740a5ad7500cd9d86c6a57a6f8b81619670697e86de4fe8eaac634f0330b"

	for _, user := range users {
		msg := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		fmt.Println(msg)
		signature1, pub, address := Sign([]byte(keyWif), []byte(msg))

		log.Printf("%x", signature1)
		signature2, pub2, address2, err := SignWIFACDSA(keyWif, []byte(msg))
		if err != nil {
			log.Println(err)
		}
		log.Printf("%x", signature2)
		signature3, pub3, address3 := Sign(key, []byte(msg))

		if err != nil {
			log.Println(err)
		}
		log.Printf("%x", signature3)
		pubKey1, err := btcec.ParsePubKey(pub)
		if err != nil {
			log.Println(err)
		}

		pubKey2, err := btcec.ParsePubKey(pub2)
		if err != nil {
			log.Println(err)
		}

		pubKey3, err := btcec.ParsePubKey(pub3)
		if err != nil {
			log.Println(err)
		}

		wantAddress := resultAddress
		wantPub := resultPubKey

		if !ecdsa.VerifyASN1(pubKey1.ToECDSA(), []byte(msg), signature1) {
			t.Errorf("got %q, wanted", signature1)
		}
		if !ecdsa.VerifyASN1(pubKey2.ToECDSA(), []byte(msg), signature2) {
			t.Errorf("got %q, wanted", signature1)
		}

		if !ecdsa.VerifyASN1(pubKey3.ToECDSA(), []byte(msg), signature3) {
			t.Errorf("got %q, wanted", signature1)
		}

		if !bytes.Equal(address, address3) {
			t.Errorf("got %q, wanted %q", address, address3)

		}
		if !bytes.Equal(address, address2) {
			t.Errorf("got %q, wanted %q", address, address2)

		}

		if !bytes.Equal(address, wantAddress) {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}
		pubString := fmt.Sprintf("%x", pub)
		if wantPub != pubString {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}

	}

}
