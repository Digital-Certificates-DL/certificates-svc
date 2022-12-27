package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	ecdsa2 "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

func Sign(key string, hashText []byte) ([]byte, []byte, []byte, error) {
	wifKey, err := btcutil.DecodeWIF(key)
	if err != nil {
		return nil, nil, nil, err
	}

	privKey := secp256k1.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
	sign := ecdsa2.SignCompact(privKey, hashText, false)

	pub, err := btcec.ParsePubKey(privKey.PubKey().SerializeUncompressed())
	if err != nil {

		return nil, nil, nil, err
	}

	return []byte(base64.RawStdEncoding.EncodeToString(sign)), pub.SerializeCompressed(), GenerateAddress(pub), err
}

//func SignWIF(key string, hashText []byte) ([]byte, []byte, []byte, error) {
//	//wifKey, err := btcutil.DecodeWIF(key)
//	//if err != nil {
//	//	return nil, nil, nil, err
//	//}
//	//log.Printf("prvate Key: %s", wifKey.String())
//	//
//	//sign, err := wifKey.PrivKey.Sign(hashText)
//	//if err != nil {
//	//	return nil, nil, nil, err
//	//}
//	privKey := secp256k1.PrivKeyFromBytes([]byte(key))
//	sign := ecdsa2.SignCompact(privKey, hashText, true)
//
//	pub, err := btcec.ParsePubKey(privKey.PubKey().SerializeUncompressed())
//	if err != nil {
//
//		return nil, nil, nil, err
//	}
//
//	return sign.Serializ0e(), wifKey.SerializePubKey(), GenerateAddress(pub), err
//}

func SignWIFACDSA(key string, hashText []byte) ([]byte, []byte, []byte, error) {
	wifKey, err := btcutil.DecodeWIF(key)
	if err != nil {
		return nil, nil, nil, err
	}

	sign, err := wifKey.PrivKey.ToECDSA().Sign(rand.Reader, hashText, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	pub, err := btcec.ParsePubKey(wifKey.SerializePubKey())
	if err != nil {

		return nil, nil, nil, err
	}

	return sign, wifKey.SerializePubKey(), GenerateAddress(pub), err
}

func GenerateAddressFromBytes(key []byte) []byte {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []byte(mainNetAddrUn.EncodeAddress())

}

func GenerateAddress(key *btcec.PublicKey) []byte {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []byte(mainNetAddrUn.EncodeAddress())

}
