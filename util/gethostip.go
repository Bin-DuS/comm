package util

import "net"
var HostIP = ""
func GetHostIP() string {
	if HostIP != "" {
		return HostIP
	}
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagLoopback != 0 {
			continue
		}
		addressList, err := inter.Addrs()
		if err != nil {
			return ""
		}
		for _, a := range addressList {
			ipNet, ok := a.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				if ip4[0] == 10 || ip4[0] == 172 || ip4[0] == 192 {
					HostIP = ip4.String()
					return HostIP
				}
			}
		}
	}
	return ""
}
