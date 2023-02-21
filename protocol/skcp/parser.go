package skcp

import (
	"errors"
	"encoding/binary"

	"github.com/LilithGames/spiracle/protocol"
)


func Parser() protocol.FuncTokenParser {
	size := 25
	return func(data []byte) (uint32, error) {
		if len(data) < size {
			return 0, errors.New("invalid skcp data")
		}
		token := binary.LittleEndian.Uint32(data[1:])
		return token, nil
	}
}
