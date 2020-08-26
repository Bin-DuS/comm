package util

import "testing"

func BenchmarkGetHostIP(b *testing.B) {
	for i := 0;i < b.N;i++ {
		GetHostIP()
	}
}
