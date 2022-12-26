package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func Sign(priv, hashText []byte) ([]byte, []byte) {

	privKey, pubKey := btcec.PrivKeyFromBytes(priv)

	signature := ecdsa.Sign(privKey, hashText)

	mainNetAddr, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	return signature.Serialize(), pubKey.SerializeCompressed()
}
