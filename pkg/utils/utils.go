package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParsePorts parses a port string like "80,443,8080-8090" into a slice of ints.
func ParsePorts(portStr string) ([]int, error) {
	var ports []int
	parts := strings.Split(portStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", rangeParts[1])
			}

			if start > end || start < 1 || end > 65535 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			port, err := strconv.Atoi(part)
			if err != nil || port < 1 || port > 65535 {
				return nil, fmt.Errorf("invalid port number: %s", part)
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

// ParseCIDR parses a CIDR notation or single IP/hostname into a list of targets.
func ParseCIDR(input string) ([]string, error) {
	if strings.Contains(input, "/") {
		_, ipnet, err := net.ParseCIDR(input)
		if err != nil {
			return nil, err
		}

		var ips []string
		ip := ipnet.IP.Mask(ipnet.Mask)
		for ; ipnet.Contains(ip); incrementIP(ip) {
			ips = append(ips, ip.String())
		}
		// Remove network and broadcast addresses for IPv4
		if ipnet.IP.To4() != nil && len(ips) > 2 {
			ips = ips[1 : len(ips)-1]
		}
		return ips, nil
	}

	return []string{input}, nil
}

func incrementIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
