package util

import "testing"

func TestMain(m *testing.M){
	InitServerName("test")
	m.Run()
}
