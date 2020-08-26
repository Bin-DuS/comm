package util

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
)

var ServerName = ""
func InitServerName(serverName string){
	ServerName = serverName
}
func CreateWorkID(url,request string)string{
	if ServerName == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%s_%s_%s",ServerName,GetTime(),url,request)))
}
func ParseWorkID(workID string) (string,error){
	b,err := base64.StdEncoding.DecodeString(workID)
	if err != nil {
		return "", errors.WithMessage(err,"ParseWorkID failed")
	}
	return string(b),nil
}
