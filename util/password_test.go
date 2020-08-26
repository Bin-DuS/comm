package util

import (
	"fmt"
	"testing"
)

func TestPassWordCheck(t *testing.T) {
	data := []struct{
		passWord string
		minLen int
		maxLen int
		upperCase bool
		lowercase bool
		number bool
		special bool
		ret bool
	}{
		{"password",8,16,true,true,true,true,false},
		{"pass",8,16,false,false,false,false,false}, //密码长度不足
		{"password",8,16,false,false,false,false,true},
		{"Password1",8,16,true,true,true,false,true},
		{"1Password",8,16,true,true,true,false,true},
	}
	for _,d := range data {
		got := PassWordCheck(d.passWord,d.minLen,d.maxLen,d.upperCase,d.lowercase,d.number,d.special)
		if got != d.ret {
			t.Fatal(fmt.Sprintf("got(%v) != d.ret(%v)",got,d.ret))
		}
	}
}
func TestPassword(t *testing.T) {
	data := []struct{
		passWord string
		want	error
	}{
		{"password",nil},
		{"",nil},
	}
	for _,d := range data {
		if PasswordAuth(d.passWord,PasswordToDB(d.passWord)) != d.want {
			t.Fatal(fmt.Sprintf("passWord(%s) want(%v)",d.passWord,d.want))
		}
	}
}
