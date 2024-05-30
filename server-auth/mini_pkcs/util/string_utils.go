package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(str string) string {
	data := []byte(str) // 要计算哈希值的数据
	hash := md5.Sum(data)
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func EncodeBase64(data []byte) string {
	// 编码字节数组为 Base64 字符串
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str
}
