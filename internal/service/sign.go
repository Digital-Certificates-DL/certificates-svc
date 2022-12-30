package service

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
	return []byte(fmt.Sprint(base64.RawStdEncoding.EncodeToString(sign), "=")), pub.SerializeUncompressed(), GenerateAddress(pub), err
}

func GenerateAddress(key *btcec.PublicKey) []byte {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []byte(mainNetAddrUn.EncodeAddress())

}
