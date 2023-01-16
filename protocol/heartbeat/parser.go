package heartbeat

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/LilithGames/spiracle/protocol"
)

func Parser() protocol.FuncTokenParser {
	return func(data []byte) (uint32, error) {
		index := bytes.IndexByte(data, ' ')
		if index < 0 {
			return 0, errors.New("miss token delimiter")
		}

		s := string(data[0:index])
		token, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("parser token %w", err)
		}
		return uint32(token), err
	}
}
