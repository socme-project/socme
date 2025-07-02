package utils

import (
	"strconv"
	"strings"
)

func IsValidIpOrDomain(input string) bool {
	return IsValidIp(input) || IsValidDomain(input)
}

func IsValidIp(ip string) bool {
	// Simple IP validation (IPv4)
	if len(ip) == 0 {
		return false
	}
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

func IsValidDomain(domain string) bool {
	// Simple domain validation
	if len(domain) == 0 {
		return false
	}
	// Check for valid characters and length
	for _, char := range domain {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != '.' &&
			char != '-' {
			return false
		}
	}
	// Check for consecutive dots or dashes
	if strings.Contains(domain, "..") || strings.Contains(domain, "--") {
		return false
	}
	return true
}

func IsValidPort(port string) bool {
	// Simple port validation
	if len(port) == 0 {
		return false
	}
	num, err := strconv.Atoi(port)
	if err != nil || num < 1 || num > 65535 {
		return false
	}
	return true
}
