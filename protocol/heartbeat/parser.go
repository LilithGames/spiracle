package heartbeat

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncParserHandler {
	return func(data []byte) (interface{}, error) {
		index := bytes.IndexByte(data, ' ')
		if index < 0 {
			return nil, errors.New("miss token delimiter")
		}

		s := string(data[1:index])
		token, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("parser token %w", err)
		}
		return uint32(token), err
	}
}
