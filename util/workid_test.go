package util

import "testing"

func TestAll(t *testing.T) {
	workID := CreateWorkID("/user/login","{\"userName\":\"Jobs\",\"password\":\"123456\"}")
	s,_:= ParseWorkID(workID)
	t.Log(s)
}
