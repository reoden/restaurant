package common

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encode 计算字符串md5散列
func MD5Encode(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}