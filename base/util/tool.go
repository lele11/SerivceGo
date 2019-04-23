package util

import (
	"math/rand"
	"strings"
)

const (
	// RandomIntToStringSource 数字 <=> 字符串 编码源数据
	RandomIntToStringSource = "E5FCDG3HQA4B1NOPIJ2RSTUV67MWX89KLYZ"
	// RandomStringSource 随机字符串源数据
	RandomStringSource = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// EncodeStringCode 数字编码为字符串
func EncodeStringCode(num uint64) string {
	code := ""
	for num > 0 {
		mod := num % 35
		num = (num - mod) / 35
		code = string(RandomIntToStringSource[mod]) + code
	}
	for len(code) < 6 {
		code = "0" + code
	}

	return code
}

// DecodeStringCode 解析字符串为数字
func DecodeStringCode(code string) uint64 {
	code = strings.Replace(code, "0", "", -1)
	var num uint64
	for i, length := 0, len(code); i < length; i++ {
		var pos uint64
		for k, v := range RandomIntToStringSource {
			if uint8(v) == code[i] {
				pos = uint64(k)
				break
			}
		}
		t := uint64(1)
		for j := 0; j < length-i-1; j++ {
			t = t * 35
		}
		num += pos * t // uint64(math.Pow(35, float64(length-i-1)))
	}
	return num
}

// RandomString 随机字符串
func RandomString(rand *rand.Rand) string {
	bytes := []byte(RandomStringSource)
	var result []byte
	for i := 0; i < 16; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
