package agent

import (
	"reflect"
	"testing"
)

func TestNormalizeAgentIPsFiltersInvalidAndUnhelpfulAddresses(t *testing.T) {
	got := normalizeAgentIPs([]string{
		"127.0.0.1",
		"192.168.1.10",
		"10.0.0.5",
		"192.168.1.10",
		"fe80::1",
		"not-an-ip",
		"0.0.0.0",
	})
	want := []string{"10.0.0.5", "192.168.1.10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected ips: got %#v, want %#v", got, want)
	}
}
