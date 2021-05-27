package multiplex

import (
	"errors"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncParserHandler {
	minSize := 1
	return func(data []byte) (interface{}, error) {
		if len(data) < minSize {
			return nil, errors.New("invalid multiplex buffer")
		}
		b := data[0]
		switch b {
		case 0x01:
			fallthrough
		case 'x':
			return b, nil
		default:
			return nil, errors.New("unknown channel")
		}
	}
}
