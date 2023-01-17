package signature

import (
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Signature struct {
	address []byte
	pubKey  *btcec.PublicKey
	key     *btcec.PrivateKey
}

func NewSignature(key string) (Signature, error) {
	wifKey, err := btcutil.DecodeWIF(key)
	if err != nil {
		return Signature{}, errors.Wrap(err, "unsupported type of key")
	}
	private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())
	pub, err := btcec.ParsePubKey(public.SerializeUncompressed())
	if err != nil {
		return Signature{}, errors.Wrap(err, "failed to parse pub key")
	}

	sign := Signature{
		key:    private,
		pubKey: pub,
	}

	address, err := sign.GenerateAddress()
	if err != nil {
		return Signature{}, errors.Wrap(err, "failed to generate address")
	}
	sign.address = address
	return sign, nil

}

func (s Signature) Sign(msg string) ([]byte, []byte, []byte, error) {
	msgWithSuf, err := s.CreateMagicMessage(msg)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to prepare msg")
	}
	messageHash := chainhash.DoubleHashB([]byte(msgWithSuf))
	sign := ecdsa.SignCompact(s.key, messageHash, false)
	return []byte(fmt.Sprint(base64.RawStdEncoding.EncodeToString(sign), "=")), s.pubKey.SerializeUncompressed(), s.address, err
}

func (s Signature) GenerateAddress() ([]byte, error) {
	mainNetAddrUn, err := btcutil.NewAddressPubKey(s.pubKey.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, errors.Wrap(err, "unsupported params")
	}
	return []byte(mainNetAddrUn.EncodeAddress()), nil

}
