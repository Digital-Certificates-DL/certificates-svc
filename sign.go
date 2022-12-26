package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func Sign(priv, hashText []byte) ([]byte, []byte, []byte) {
	privKey, pubKey := btcec.PrivKeyFromBytes(priv)
	signature := ecdsa.Sign(privKey, hashText)
	return signature.Serialize(), pubKey.SerializeCompressed(), GenerateAddress(pubKey)
}

func GenerateAddress(key *btcec.PublicKey) []byte {
	mainNetAddr, err := btcutil.NewAddressPubKey(key.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	mainNetAddr.SetFormat(btcutil.PKFCompressed)
	return []byte(mainNetAddr.EncodeAddress())

}
