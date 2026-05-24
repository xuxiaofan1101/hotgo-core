package agent

import (
	"net"
	"sort"
	"strings"
)

func localAgentIPs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	items := make([]string, 0)
	for _, item := range interfaces {
		if item.Flags&net.FlagUp == 0 || item.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := item.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				items = append(items, v.IP.String())
			case *net.IPAddr:
				items = append(items, v.IP.String())
			}
		}
	}
	return normalizeAgentIPs(items)
}

func normalizeAgentIPs(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		ip := net.ParseIP(strings.TrimSpace(item))
		if ip == nil || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			continue
		}
		value := ip.String()
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
