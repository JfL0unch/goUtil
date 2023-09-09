package utils

import (
	"crypto/md5"
	"crypto/sha1"
)

// Hash 类比java: DigestUtils.sha1Hex
func Hash(toSign []byte) []byte {
	hasher := sha1.New()
	hasher.Write(toSign)
	return hasher.Sum(nil)
}

func Md5(toSign []byte) []byte {
	m5 := md5.New()
	m5.Write(toSign)
	return m5.Sum(nil)
}
