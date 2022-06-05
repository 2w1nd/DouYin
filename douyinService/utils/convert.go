package utils

import (
	"strconv"
	"strings"
)

func String2Uint64(str string) uint64 {
	atoi, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return atoi
}

func String2Int64(str string) int64 {
	atoi, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return atoi
}

func SplitString(str, c string) (res string) {
	idx := strings.LastIndex(str, c)
	res = str[idx+1:]
	return
}
