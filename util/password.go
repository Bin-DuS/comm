package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

//使用 8-32个字符（至少包含大小写字母和数字）；
/*
	passWord  密码字符串
	minLen	  最短长度 至少大于1
	maxLen    最大长度 0 为不限制
	upperCase 是否有大写字母
	lowercase 是否有小写字母
	number    是否有数字
	special   是否有特殊字符
*/
func PassWordCheck(passWord string,minLen,maxLen int,upperCase,lowercase,number,special bool) bool {
	if len(passWord) < minLen || minLen <= 0 {
		return false
	}
	if len(passWord) > maxLen && maxLen != 0 { // maxLen = 0 不限制最大长度
		return false
	}
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range passWord {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if upperCase == true && hasUpperCase != true {
		return false
	}
	if lowercase == true && hasLowercase != true {
		return false
	}
	if number == true && hasNumber != true {
		return false
	}
	if special == true && hasSpecial != true {
		return false
	}
	return true
}

func PasswordToDB(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	return string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即可
}
func PasswordAuth(password,passwordFromDB string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password)) //验证（对比）
}
