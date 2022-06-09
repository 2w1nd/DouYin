package codes

// 状态码定义
const (
	SUCCESS = 0
	ERROR   = 1

	Followed = 101
	Follow   = 102
	FOCUS    = 103
	NoFOCUS  = 104

	BITMAPLIKE    int = 201
	BITMAPUNLIKE  int = 202
	ALREADYEXIST  int = 203
	ALREADYDELETE int = 204

	RedisNotFound int = 404
	RedisFound    int = 405
)
