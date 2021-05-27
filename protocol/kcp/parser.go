package kcp

import (
	"encoding/binary"
	"errors"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncParserHandler {
	ikcpOverhead := 24
	return func(data []byte) (interface{}, error) {
		if len(data) < ikcpOverhead {
			return nil, errors.New("invalid kcp data")
		}
		token := binary.LittleEndian.Uint32(data)
		return token, nil
	}
}
