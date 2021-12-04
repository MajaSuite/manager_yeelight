package device

import "net"

func CheckLan(host string) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, YeeLightPort), Timeout); err == nil {
		defer conn.Close()
		return true
	}
	return false
}

