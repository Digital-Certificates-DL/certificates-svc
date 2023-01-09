package signature

import (
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

type Signature struct {
	msg string
	key string
}

func NewSignature(msg, key string) Signature {
	return Signature{
		msg: msg,
		key: key,
	}
}

func (s Signature) Sign() ([]byte, []byte, []byte, error) {
	wifKey, err := btcutil.DecodeWIF(s.key)
	if err != nil {
		return nil, nil, nil, err
	}

	msgWithSuf := s.CreateMagicMessage(s.msg)

	private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())

	pub, err := btcec.ParsePubKey(public.SerializeUncompressed())
	if err != nil {
		return nil, nil, nil, err
	}
	messageHash := chainhash.DoubleHashB([]byte(msgWithSuf))

	sign := ecdsa.SignCompact(private, messageHash, false)
	return []byte(fmt.Sprint(base64.RawStdEncoding.EncodeToString(sign), "=")), pub.SerializeUncompressed(), s.GenerateAddress(pub), err
}

func (s Signature) GenerateAddress(key *btcec.PublicKey) []byte {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return []byte(mainNetAddrUn.EncodeAddress())

}
