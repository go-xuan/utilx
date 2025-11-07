package randx

import (
	"hash/crc32"
	"math/rand"

	"github.com/google/uuid"
)

// 字典配置
const (
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers      = "1234567890"
	special      = "~!@#$%^&*()-+_=:;,|./?"
	lowerChar    = "abcdefghijklmnopqrstuvwxyz1234567890"
	allChar      = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var _rand *rand.Rand

func NewRand() *rand.Rand {
	if _rand == nil {
		_rand = rand.New(rand.NewSource(seed()))
	}
	return _rand
}

// 随机种子
func seed() int64 {
	return int64(crc32.ChecksumIEEE([]byte(uuid.NewString())))
}
