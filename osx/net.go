package osx

import (
	"net"
	"strconv"
	"strings"
)

// GetLocalIP 获取本地IP
func GetLocalIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		return getIPv4(addrs)
	}
	return ""
}

// GetWLANIP 获取当前机器的WLAN Host
func GetWLANIP() string {
	if interfaces, err := net.Interfaces(); err == nil {
		for _, in := range interfaces {
			if (in.Flags&net.FlagUp) != 0 && in.Name == "WLAN" {
				if addrs, e := in.Addrs(); e == nil {
					if ip := getIPv4(addrs); ip != "" {
						return ip
					}
				}
			}
		}
	}
	return ""
}

// getIPv4 从地址列表中查找有效的 IPv4 地址
func getIPv4(addrs []net.Addr) string {
	for _, addr := range addrs {
		if n, ok := addr.(*net.IPNet); ok && !n.IP.IsLoopback() {
			if ip := n.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}

// CheckIpMatch 检测IP是否匹配规则
func CheckIpMatch(rules []string, ip string) bool {
	if len(rules) == 0 || ip == "" {
		return false
	}

	for _, rule := range rules {
		if strings.Contains(rule, `-`) {
			ruleStart, ruleEnd, found := strings.Cut(rule, `-`)
			if !found {
				continue
			}
			prefix, num := SplitIpByLastPoint(ip)
			startPrefix, minNum := SplitIpByLastPoint(ruleStart)
			endPrefix, maxNum := SplitIpByLastPoint(ruleEnd)
			if prefix == startPrefix && prefix == endPrefix && num >= minNum && num <= maxNum {
				return true
			}
		} else if ipMatch(rule, ip) {
			return true
		}
	}
	return false
}

// ipMatch 检测IP是否匹配规则
func ipMatch(rule, ip string) bool {
	ruleParts := strings.Split(rule, ".")
	ipParts := strings.Split(ip, ".")
	match := true
	for i := 0; i < len(ruleParts); i++ {
		if ruleParts[i] != "*" && ruleParts[i] != ipParts[i] {
			match = false
			break
		}
	}
	return match
}

// SplitIpByLastPoint 将IP以最后一个.拆分
func SplitIpByLastPoint(ip string) (string, int) {
	lastIndex := strings.LastIndex(ip, ".")
	if lastIndex == -1 {
		return "", 0
	}

	prefix := ip[:lastIndex]
	sufStr := ip[lastIndex+1:]
	suf, err := strconv.Atoi(sufStr)
	if err != nil {
		return prefix, 0
	}
	return prefix, suf
}
