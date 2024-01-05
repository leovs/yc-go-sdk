package eureka_client

import (
	"net"
	"os"
)

// getLocalIP 获取本地ip
func getLocalIP() string {

	// 支持通过环境变量设置IP地址
	address := os.Getenv("EUREKA_INSTANCE_IP-ADDRESS")
	if address != "" {
		return address
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, interFace := range netInterfaces {
		// 只获取已上线网卡
		if (interFace.Flags & net.FlagUp) != 0 {
			addrs, _ := interFace.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}
