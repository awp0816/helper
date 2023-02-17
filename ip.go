package helper

import (
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
)

func GetLocalIp() string {
	ip := ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() && !ipnet.IP.IsLinkLocalUnicast() && !ipnet.IP.IsLinkLocalMulticast() && ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func IsPortInUse(host string, port int64) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", port)), time.Second*1)
	if err == nil {
		_ = conn.Close()
		return true
	}

	return false
}

func Ping(ip string) (bool, string) {
	pinger := ping.New(ip)
	pinger.Count = 4
	pinger.Timeout = time.Duration(2000) * time.Millisecond
	pinger.Interval = time.Duration(500) * time.Millisecond
	pinger.SetPrivileged(true)
	if err := pinger.Run(); err != nil {
		return false, err.Error()
	}
	if pinger.Statistics().PacketsRecv >= 1 {
		return true, fmt.Sprintf("设备[%s]可以访问", ip)
	}
	return false, fmt.Sprintf("设备[%s]网络可能不太稳定", ip)
}
