package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(content string) string {
	h := md5.New()
	h.Write([]byte(content)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}
