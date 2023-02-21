package multiplex

import (
	"errors"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncMultiplexParser {
	size := 1
	return func(data []byte) (byte, error) {
		if len(data) < size {
			return 0, errors.New("invalid multiplex buffer")
		}
		b := data[0]
		switch b {
		case 0x01:
			fallthrough
		case 'e':
			fallthrough
		case 's':
			fallthrough
		case 'x':
			return b, nil
		default:
			return 0, errors.New("unknown multiplexing channel")
		}
	}
}
