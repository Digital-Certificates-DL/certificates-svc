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
		return nil, nil, nil, errors.Wrap(err, "unsupported type of key")
	}

	msgWithSuf, err := s.CreateMagicMessage(s.msg)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to prepare msg")
	}

	private, public := btcec.PrivKeyFromBytes(wifKey.PrivKey.Serialize())

	pub, err := btcec.ParsePubKey(public.SerializeUncompressed())
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to parse pub key")
	}
	messageHash := chainhash.DoubleHashB([]byte(msgWithSuf))

	sign := ecdsa.SignCompact(private, messageHash, false)
	address, err := s.GenerateAddress(pub)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to generate address")
	}
	return []byte(fmt.Sprint(base64.RawStdEncoding.EncodeToString(sign), "=")), pub.SerializeUncompressed(), address, err
}

func (s Signature) GenerateAddress(key *btcec.PublicKey) ([]byte, error) {

	mainNetAddrUn, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return nil, errors.Wrap(err, "unsupported params")
	}
	return []byte(mainNetAddrUn.EncodeAddress()), nil

}
