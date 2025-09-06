package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"
)

// 初始化随机数种子
func init() {
	mathRand.Seed(time.Now().UnixNano())
}

// GenerateRandomCode 生成指定长度的数字验证码
// length: 验证码长度
func GenerateRandomCode(length int) string {
	if length <= 0 {
		length = 6 // 默认6位
	}
	return fmt.Sprintf("%0*d", length, mathRand.Intn(intPow(10, length)))
}

// GenerateSecureRandomCode 生成加密安全的随机数字验证码
// length: 验证码长度
func GenerateSecureRandomCode(length int) (string, error) {
	if length <= 0 {
		length = 6 // 默认6位
	}

	max := big.NewInt(int64(intPow(10, length)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, n.Int64()), nil
}

// GenerateRandomString 生成指定长度的随机字符串
// length: 字符串长度
// charset: 字符集，如果为空则使用默认字符集
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	if length <= 0 {
		length = 8 // 默认8位
	}

	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charset[mathRand.Intn(len(charset))])
	}
	return result.String()
}

// GenerateSecureRandomString 生成加密安全的随机字符串
// length: 字符串长度
// charset: 字符集，如果为空则使用默认字符集
func GenerateSecureRandomString(length int, charset string) (string, error) {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	if length <= 0 {
		length = 8 // 默认8位
	}

	var result strings.Builder
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result.WriteByte(charset[num.Int64()])
	}
	return result.String(), nil
}

// GenerateRandomInt 生成指定范围内的随机整数 [min, max)
// min: 最小值（包含）
// max: 最大值（不包含）
func GenerateRandomInt(min, max int) int {
	if min >= max {
		return min
	}
	return mathRand.Intn(max-min) + min
}

// GenerateSecureRandomInt 生成加密安全的指定范围内的随机整数 [min, max)
// min: 最小值（包含）
// max: 最大值（不包含）
func GenerateSecureRandomInt(min, max int) (int, error) {
	if min >= max {
		return min, nil
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}

// GenerateRandomFloat64 生成指定范围内的随机浮点数 [min, max)
// min: 最小值（包含）
// max: 最大值（不包含）
func GenerateRandomFloat64(min, max float64) float64 {
	if min >= max {
		return min
	}
	return mathRand.Float64()*(max-min) + min
}

// GenerateRandomBool 生成随机布尔值
func GenerateRandomBool() bool {
	return mathRand.Intn(2) == 1
}

// GenerateRandomBytes 生成指定长度的随机字节数组
// length: 字节数组长度
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateUUID 生成简单的UUID（不完全符合RFC标准，仅用于简单场景）
func GenerateUUID() string {
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		// 如果加密随机数生成失败，使用数学随机数作为备选
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			mathRand.Uint32(),
			mathRand.Uint32()&0xffff,
			mathRand.Uint32()&0xffff,
			mathRand.Uint32()&0xffff,
			mathRand.Uint64()&0xffffffffffff)
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

// 预定义的字符集常量
const (
	// DigitCharset 数字字符集
	DigitCharset = "0123456789"
	// LowerCharset 小写字母字符集
	LowerCharset = "abcdefghijklmnopqrstuvwxyz"
	// UpperCharset 大写字母字符集
	UpperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// AlphaCharset 字母字符集
	AlphaCharset = LowerCharset + UpperCharset
	// AlphaNumCharset 字母数字字符集
	AlphaNumCharset = AlphaCharset + DigitCharset
	// SpecialCharset 特殊字符集
	SpecialCharset = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	// AllCharset 所有字符集
	AllCharset = AlphaNumCharset + SpecialCharset
)

// 便捷函数

// GenerateDigitString 生成指定长度的数字字符串
func GenerateDigitString(length int) string {
	return GenerateRandomString(length, DigitCharset)
}

// GenerateAlphaString 生成指定长度的字母字符串
func GenerateAlphaString(length int) string {
	return GenerateRandomString(length, AlphaCharset)
}

// GenerateAlphaNumString 生成指定长度的字母数字字符串
func GenerateAlphaNumString(length int) string {
	return GenerateRandomString(length, AlphaNumCharset)
}

// GeneratePassword 生成指定长度的密码（包含字母、数字和特殊字符）
func GeneratePassword(length int) string {
	if length < 4 {
		length = 8 // 最少8位
	}

	// 确保密码包含各种类型的字符
	var password strings.Builder
	
	// 至少包含一个小写字母
	password.WriteByte(LowerCharset[mathRand.Intn(len(LowerCharset))])
	// 至少包含一个大写字母
	password.WriteByte(UpperCharset[mathRand.Intn(len(UpperCharset))])
	// 至少包含一个数字
	password.WriteByte(DigitCharset[mathRand.Intn(len(DigitCharset))])
	// 至少包含一个特殊字符
	password.WriteByte(SpecialCharset[mathRand.Intn(len(SpecialCharset))])

	// 填充剩余长度
	for i := 4; i < length; i++ {
		password.WriteByte(AllCharset[mathRand.Intn(len(AllCharset))])
	}

	// 打乱字符顺序
	passwordBytes := []byte(password.String())
	mathRand.Shuffle(len(passwordBytes), func(i, j int) {
		passwordBytes[i], passwordBytes[j] = passwordBytes[j], passwordBytes[i]
	})

	return string(passwordBytes)
}

// 辅助函数：计算整数幂
func intPow(base, exp int) int {
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}