package tools

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// IPInRanges checks if IP is within specified ranges
// ipRangesJSON: JSON array of IP ranges, e.g.: ["192.168.1.0/24", "10.0.0.0/8", "192.168.2.100"]
func IPInRanges(clientIP string, ipRangesJSON string) (bool, error) {
	if ipRangesJSON == "" {
		return false, fmt.Errorf("IP ranges list is empty")
	}

	var ipRanges []string
	if err := json.Unmarshal([]byte(ipRangesJSON), &ipRanges); err != nil {
		return false, fmt.Errorf("failed to parse IP ranges JSON: %v", err)
	}

	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", clientIP)
	}

	for _, ipRange := range ipRanges {
		ipRange = strings.TrimSpace(ipRange)
		if ipRange == "" {
			continue
		}

		if strings.Contains(ipRange, "/") {
			_, network, err := net.ParseCIDR(ipRange)
			if err != nil {
				continue
			}
			if network.Contains(ip) {
				return true, nil
			}
		} else {
			rangeIP := net.ParseIP(ipRange)
			if rangeIP != nil && rangeIP.Equal(ip) {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetClientIP extracts client IP from request headers
func GetClientIP(remoteAddr string, xForwardedFor string, xRealIP string) string {
	if xRealIP != "" {
		ip := strings.TrimSpace(strings.Split(xRealIP, ",")[0])
		if ip != "" {
			return ip
		}
	}

	if xForwardedFor != "" {
		ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		if ip != "" {
			return ip
		}
	}

	if remoteAddr != "" {
		ip := strings.Split(remoteAddr, ":")[0]
		return ip
	}

	return ""
}

// ValidateIPRange validates IP range format
func ValidateIPRange(ipRange string) error {
	ipRange = strings.TrimSpace(ipRange)
	if ipRange == "" {
		return fmt.Errorf("IP range cannot be empty")
	}

	if strings.Contains(ipRange, "/") {
		_, _, err := net.ParseCIDR(ipRange)
		if err != nil {
			return fmt.Errorf("invalid CIDR format: %v", err)
		}
	} else {
		ip := net.ParseIP(ipRange)
		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", ipRange)
		}
	}

	return nil
}
