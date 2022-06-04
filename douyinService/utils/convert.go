package utils

import "strconv"

func String2Uint64(str string) uint64 {
	atoi, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return atoi
}
