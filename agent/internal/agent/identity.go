package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var (
	readStableMachineId      = defaultReadStableMachineId
	readPhysicalMacAddresses = defaultReadPhysicalMacAddresses
)

func deriveStableAgentId(machineId string, macAddresses []string) string {
	machineId = strings.TrimSpace(machineId)
	macs := normalizeMacAddresses(macAddresses)
	if machineId == "" && len(macs) == 0 {
		return ""
	}
	parts := []string{runtime.GOOS}
	if machineId != "" {
		parts = append(parts, machineId)
	}
	parts = append(parts, macs...)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return "agent-" + hex.EncodeToString(sum[:])[:24]
}

func defaultReadStableMachineId() string {
	switch runtime.GOOS {
	case "linux":
		return readFirstNonEmptyFile("/etc/machine-id", "/var/lib/dbus/machine-id")
	case "darwin":
		return readDarwinPlatformUUID()
	case "windows":
		return readWindowsMachineGUID()
	default:
		return ""
	}
}

func readFirstNonEmptyFile(paths ...string) string {
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if value := strings.TrimSpace(string(content)); value != "" {
			return value
		}
	}
	return ""
}

func readDarwinPlatformUUID() string {
	output, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return ""
	}
	matches := regexp.MustCompile(`"IOPlatformUUID"\s*=\s*"([^"]+)"`).FindStringSubmatch(string(output))
	if len(matches) != 2 {
		return ""
	}
	return strings.TrimSpace(matches[1])
}

func readWindowsMachineGUID() string {
	output, err := exec.Command("reg", "query", `HKLM\SOFTWARE\Microsoft\Cryptography`, "/v", "MachineGuid").Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(output))
	for i, field := range fields {
		if strings.EqualFold(field, "MachineGuid") && i+2 < len(fields) {
			return strings.TrimSpace(fields[i+2])
		}
	}
	return ""
}

func defaultReadPhysicalMacAddresses() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	items := make([]string, 0)
	for _, item := range interfaces {
		if item.Flags&net.FlagLoopback != 0 || len(item.HardwareAddr) == 0 || isVirtualInterfaceName(item.Name) {
			continue
		}
		items = append(items, item.HardwareAddr.String())
	}
	return normalizeMacAddresses(items)
}

func normalizeMacAddresses(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		value := strings.ToLower(strings.TrimSpace(item))
		if value == "" {
			continue
		}
		mac, err := net.ParseMAC(value)
		if err != nil || len(mac) == 0 {
			continue
		}
		value = mac.String()
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func isVirtualInterfaceName(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return false
	}
	virtualPrefixes := []string{
		"awdl", "br-", "bridge", "docker", "llw", "tap", "tailscale", "tun", "utun", "veth", "virbr", "vmnet", "wg", "zt",
	}
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}
