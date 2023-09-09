package utils

import (
	"crypto/rand"
	"math/big"
	pseudorand "math/rand"
	url2 "net/url"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandN 随机数
// pseudo: true 伪随机, false 真随机
func RandN(n int, pseudo bool) int {
	if pseudo {
		return getPseudoRand(n)
	} else {
		return int(getRealRand(int64(n)).Int64())
	}
}

// RandString 随机字符串
// pseudo: true 伪随机, false 真随机
func RandString(n int, pseudo bool) string {

	b := make([]byte, n)
	for i := range b {
		if pseudo {
			c := getPseudoRand(n)
			b[i] = letterBytes[c]
		} else {
			c := int(getRealRand(int64(n)).Int64())
			b[i] = letterBytes[c]
		}
	}
	return string(b)
}

func getPseudoRand(n int) int {
	return pseudorand.Intn(n)
}

func getRealRand(n int64) *big.Int {
	a, _ := rand.Int(rand.Reader, big.NewInt(n))
	return a
}

// JavaScriptEncodeURI
// js encodeURI
// golang版本
// https://blog.csdn.net/m0_46309087/article/details/119839122
func JavaScriptEncodeURI(s string) string {

	s = url2.QueryEscape(s)
	metaStr := map[string]string{
		"%3B": ";",
		"%2C": ",",
		"%2F": "/",
		"%3F": "?",
		"%3A": ":",
		"%40": "@",
		"%26": "&",
		"%3D": "=",
		"%2B": "+",
		"%24": "$",
		"%23": "#",
	}
	for from, to := range metaStr {
		s = strings.Replace(s, from, to, -1)
	}
	return s
}
