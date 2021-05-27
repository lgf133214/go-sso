package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(data string) string {
	tmp := md5.Sum([]byte(data))
	tmp2 := make([]byte, 32)
	hex.Encode(tmp2, tmp[:])
	return string(tmp2)
}
