package repos

import "strconv"

func str(token TToken) string {
	return strconv.FormatUint(uint64(token), 16)
}
