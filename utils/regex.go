package utils

import (
	"regexp"
	"strings"
)

// RegexPatterns 正则表达式模式常量
type RegexPatterns struct{}

// 正则表达式常量定义
const (
	// PhoneRegex 手机号正则
	PhoneRegex = `^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`
	// EmailRegex 邮箱正则
	EmailRegex = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	// PasswordRegex 密码正则。4~32位的字母、数字、下划线
	PasswordRegex = `^\w{4,32}$`
	// VerifyCodeRegex 验证码正则, 6位数字或字母
	VerifyCodeRegex = `^[a-zA-Z\d]{6}$`
)

// 预编译正则表达式以提高性能
var (
	phoneRegexp      = regexp.MustCompile(PhoneRegex)
	emailRegexp      = regexp.MustCompile(EmailRegex)
	passwordRegexp   = regexp.MustCompile(PasswordRegex)
	verifyCodeRegexp = regexp.MustCompile(VerifyCodeRegex)
)

// IsPhoneInvalid 是否是无效手机格式
// phone: 要校验的手机号
// 返回 true: 不符合，false: 符合
func IsPhoneInvalid(phone string) bool {
	return mismatch(phone, phoneRegexp)
}

// IsEmailInvalid 是否是无效邮箱格式
// email: 要校验的邮箱
// 返回 true: 不符合，false: 符合
func IsEmailInvalid(email string) bool {
	return mismatch(email, emailRegexp)
}

// IsPasswordInvalid 是否是无效密码格式
// password: 要校验的密码
// 返回 true: 不符合，false: 符合
func IsPasswordInvalid(password string) bool {
	return mismatch(password, passwordRegexp)
}

// IsCodeInvalid 是否是无效验证码格式
// code: 要校验的验证码
// 返回 true: 不符合，false: 符合
func IsCodeInvalid(code string) bool {
	return mismatch(code, verifyCodeRegexp)
}

// IsPhoneValid 是否是有效手机格式
// phone: 要校验的手机号
// 返回 true: 符合，false: 不符合
func IsPhoneValid(phone string) bool {
	return !IsPhoneInvalid(phone)
}

// IsEmailValid 是否是有效邮箱格式
// email: 要校验的邮箱
// 返回 true: 符合，false: 不符合
func IsEmailValid(email string) bool {
	return !IsEmailInvalid(email)
}

// IsPasswordValid 是否是有效密码格式
// password: 要校验的密码
// 返回 true: 符合，false: 不符合
func IsPasswordValid(password string) bool {
	return !IsPasswordInvalid(password)
}

// IsCodeValid 是否是有效验证码格式
// code: 要校验的验证码
// 返回 true: 符合，false: 不符合
func IsCodeValid(code string) bool {
	return !IsCodeInvalid(code)
}

// mismatch 校验是否不符合正则格式
// str: 要校验的字符串
// regex: 预编译的正则表达式
// 返回 true: 不符合，false: 符合
func mismatch(str string, regex *regexp.Regexp) bool {
	if strings.TrimSpace(str) == "" {
		return true
	}
	return !regex.MatchString(str)
}