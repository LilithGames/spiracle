package kcp

import (
	"encoding/binary"
	"errors"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncTokenParser {
	size := 24
	return func(data []byte) (uint32, error) {
		if len(data) < size {
			return 0, errors.New("invalid kcp data")
		}
		token := binary.LittleEndian.Uint32(data)
		return token, nil
	}
}
