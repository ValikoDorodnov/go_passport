package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

type Hasher struct {
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h Hasher) GetMD5Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}
