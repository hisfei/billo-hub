package helper

import (
	"fmt"
	"net"
)

// GetLocalIP returns the first non-loopback IPv4 address found.
// If no valid IPv4 address is found, it returns an empty string.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("failed to get interface addresses: %v", err)
		return ""
	}

	for _, address := range addrs {
		// Check the address type and ensure it's not a loopback address.
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// Ensure it's an IPv4 address.
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	fmt.Println("no non-loopback IPv4 address found")

	return ""
}

// GetLocalIPs returns a slice of all non-loopback IPv4 addresses.
func GetLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("failed to get interface addresses: %v", err)

		return nil
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	if len(ips) == 0 {
		fmt.Println("no non-loopback IPv4 address found")

		return nil
	}

	return ips
}
