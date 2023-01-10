package signature

import (
	"bytes"
	"github.com/btcsuite/btcd/wire"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const varIntProtoVer uint32 = 0

const magicMessage = "\x18Bitcoin Signed Message:\n"

func (s Signature) CreateMagicMessage(message string) (string, error) {
	buffer := bytes.Buffer{}
	buffer.Grow(wire.VarIntSerializeSize(uint64(len(message))))
	if err := wire.WriteVarInt(&buffer, varIntProtoVer, uint64(len(message))); err != nil {
		return "", errors.Wrap(err, "failed to decide or insert size of msg")
	}
	return magicMessage + buffer.String() + message, nil
}
