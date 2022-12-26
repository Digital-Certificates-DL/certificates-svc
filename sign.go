package main

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func Sign(priv, hashText []byte) ([]byte, []byte) {

	privKey, pubKey := btcec.PrivKeyFromBytes(priv)

	signature := ecdsa.Sign(privKey, hashText)

	return signature.Serialize(), pubKey.SerializeCompressed()
}
