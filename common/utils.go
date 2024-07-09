package common

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"time"
	"unicode"
)

type Option func(map[string]interface{}) error

func WithConfig(key string, value interface{}) Option {
	return func(config map[string]interface{}) error {
		config[key] = value
		return nil
	}
}

func GetConfig(opts ...Option) map[string]interface{} {
	config := map[string]interface{}{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func Mask(content string) string {
	if isDigit(content) && len(content) > 10 {
		return content[:3] + "****" + content[len(content)-4:]
	}
	return content
}

func isDigit(content string) bool {
	match, _ := regexp.MatchString("^[0-9]+$", content)
	return match
}

func BytesToSize(bytes float64, count int) string {
	bytes /= 1024.0 * 1024.0
	if count == 2 {
		return fmt.Sprintf("%.2f %s", bytes, "MB")
	}
	return fmt.Sprintf("%.2f %s", bytes, "MB")
}

// 生成随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成随机数字串
func RandDigit(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成随机字符串
func RandUpperString(length int) string {
	str := "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// b64 to bytes
func B64ToBytes(b64 string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func BytesToB64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

var (
	SplitCount        = 600
	AudioSplitCount   = 250
	AudioEnSplitCount = 1500
)

// 公众号切分文本长度
func GetSplitCount(contentRune []rune) int {
	isHan := false
	for _, c := range contentRune {
		if unicode.Is(unicode.Han, c) {
			isHan = true
			break
		}
	}

	splitCount := AudioEnSplitCount
	if isHan {
		splitCount = AudioSplitCount
	}
	return splitCount
}

func StringPtr(s string) *string {
	return &s
}

func UInt64Ptr(i uint64) *uint64 {
	return &i
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
