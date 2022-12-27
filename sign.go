package main

import (
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

func Sign(key string, hashText string) ([]byte, []byte, []byte, error) {
	wifKey, err := btcutil.DecodeWIF(key)
	if err != nil {
		return nil, nil, nil, err
	}

	msg := CreateMagicMessage(hashText)

	private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())

	pub, err := btcec.ParsePubKey(public.SerializeUncompressed())
	if err != nil {
		return nil, nil, nil, err
	}
	messageHash := chainhash.DoubleHashB([]byte(msg))

	sign := ecdsa.SignCompact(private, messageHash, false)
	return []byte(base64.RawStdEncoding.EncodeToString(sign)), pub.SerializeUncompressed(), GenerateAddress(pub), err
}

//func SignCrypto(key string, hashText string) {
//	wifKey, err := btcutil.DecodeWIF(key)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	//msg := fmt.Sprintf("\x18Bitcoin Signed Message:\n%d%s", len(hashText), hashText)
//	msg := CreateMagicMessage(hashText)
//
//	messageHash := chainhash.DoubleHashB([]byte(msg))
//
//	sign, err := crypto.Sign(messageHash, wifKey.PrivKey.ToECDSA())
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	s := base64.RawStdEncoding.EncodeToString(sign)
//	log.Println("SIGN: ", s)
//
//}

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

func GenerateAddress(key *btcec.PublicKey) []byte {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []byte(mainNetAddrUn.EncodeAddress())

}

//func Sign(key string, hashText []byte) ([]byte, []byte, []byte, error) {
//	wifKey, err := btcutil.DecodeWIF(key)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	msg := fmt.Sprintf("\x19Bitcoin Signed Message:\n%d%s", len(hashText), hashText)
//
//	privKey := secp256k1.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
//	sign := ecdsa.SignCompact(privKey, []byte(msg), false)
//
//	//private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
//	pub, err := btcec.ParsePubKey(privKey.PubKey().SerializeUncompressed())
//	if err != nil {
//
//		return nil, nil, nil, err
//	}
//
//	return []byte(base64.RawStdEncoding.EncodeToString(sign)), pub.SerializeUnompressed(), GenerateAddress(pub), err
//}

//
//func Sign(key string, hashText string) ([]byte, []byte, []byte, error) {
//	wifKey, err := btcutil.DecodeWIF(key)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//
//	msg := fmt.Sprintf("\x18Bitcoin Signed Message:\n%d%s", len(hashText), hashText)
//
//	//privKey := secp256k1.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
//	log.Printf("%x", wifKey.PrivKey.Serialize())
//	private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
//	log.Printf("%x", private.Serialize())
//	pub, err := btcec.ParsePubKey(public.SerializeUncompressed())
//	if err != nil {
//
//		return nil, nil, nil, err
//	}
//	sum := sha256.Sum256([]byte(msg))
//	sign := ecdsa.Sign(private, sum[:])
//
//	fmt.Println(sign.Verify(sum[:], pub))
//	fmt.Println(sign.Serialize())
//	fmt.Println(hex.EncodeToString(sign.Serialize()))
//	//sign := ecdsa.SignCompact(private, sum[:], false)
//	return []byte(base64.RawStdEncoding.EncodeToString(sign.Serialize())), pub.SerializeUncompressed(), GenerateAddress(pub), err
//}
