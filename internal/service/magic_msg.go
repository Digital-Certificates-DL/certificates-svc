package service

import (
	"bytes"
	"github.com/btcsuite/btcd/wire"
)

const varIntProtoVer uint32 = 0

const magicMessage = "\x18Bitcoin Signed Message:\n"

func CreateMagicMessage(message string) string {
	buffer := bytes.Buffer{}
	buffer.Grow(wire.VarIntSerializeSize(uint64(len(message))))

	// If we cannot write the VarInt, just panic since that should never happen
	if err := wire.WriteVarInt(&buffer, varIntProtoVer, uint64(len(message))); err != nil {
		panic(err)
	}

	return magicMessage + buffer.String() + message
}
