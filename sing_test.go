package main

import (
	"bytes"
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

	resultAddress := []byte("1KEXKjW8R2dLUyjY4MLLi49FM1T1tyH1rx")
	keyWif := "5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ" //KwYitv9awg9xGVHJJ2SZVdE7WXYuXp9mJvSiz8XzyiNeu34WDsXo //5HtbQzoofDVmEvzAMd89iq5nRpFUm8xudgkSNZ4rPN9z5AWcGEJ //09c7d74ee3ee97c9101e308e3dcec263556751b1e0b78e788add115741f61786
	//key := "09c7d74ee3ee97c9101e308e3dcec263556751b1e0b78e788add115741f61786"
	resultPubKey := "02a423740a5ad7500cd9d86c6a57a6f8b81619670697e86de4fe8eaac634f0330b"

	for _, user := range users {
		msg := fmt.Sprintf("%s %s %s", user.Date, user.Participant, user.CourseTitle)
		fmt.Println(msg)
		signature1, pub, address, err := Sign(keyWif, msg)
		if err != nil {
			log.Println(err)
		}

		log.Printf("%x", signature1)
		SignCrypto(keyWif, msg)

		log.Printf("%x", signature1)

		//pubKey1, err := btcec.ParsePubKey(pub)
		//if err != nil {
		//	log.Println(err)
		//}

		wantAddress := resultAddress
		wantPub := resultPubKey

		//sign, err := ecdsa2.ParseDERSignature(signature1)
		//if err != nil {
		//	log.Println(err)
		//}

		//if sign.Verify([]byte(msg), pubKey1) {
		//	t.Errorf("got %q, wanted", signature1)
		//}

		if !bytes.Equal(address, wantAddress) {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}

		pubString := fmt.Sprintf("%x", pub)
		if wantPub != pubString {
			t.Errorf("got %q, wanted %q", address, wantAddress)

		}

	}

}

//func ParseCompact(signature []byte, curve *btcec.KoblitzCurve) (*ecdsa2.Signature, error) {
//	bitLen := (curve.BitSize + 7) / 8
//	if len(signature) != 1+bitLen*2 {
//		return nil, errors.New("invalid compact signature size")
//	}
//	ecdsa2.NewSignature(new(big.Int).SetBytes(signature[1:bitLen+1]), new(big.Int).SetBytes(signature[bitLen+1:]))
//	return &ecdsa2.Signature{
//		R: new(big.Int).SetBytes(signature[1 : bitLen+1]),
//		S: new(big.Int).SetBytes(signature[bitLen+1:]),
//	}, nil
//}
