package sys

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestDataAgentOnlineStatusUsesLastSeenAt(t *testing.T) {
	now := time.Date(2026, 5, 24, 20, 0, 0, 0, time.Local)

	online, label := dataAgentOnlineStatus(now.Add(-90*time.Second), now)
	if !online || label != "online" {
		t.Fatalf("expected recent heartbeat to be online, got online=%v label=%q", online, label)
	}

	online, label = dataAgentOnlineStatus(now.Add(-3*time.Minute), now)
	if online || label != "offline" {
		t.Fatalf("expected stale heartbeat to be offline, got online=%v label=%q", online, label)
	}
}

func TestDataAgentJsonStringsHandlesNilAndArrays(t *testing.T) {
	if got := dataAgentJsonStrings(nil); len(got) != 0 {
		t.Fatalf("expected nil json to return empty list, got %#v", got)
	}

	got := dataAgentJsonStrings(gjson.New([]string{"127.0.0.1", "10.0.0.1"}))
	if strings.Join(got, ",") != "127.0.0.1,10.0.0.1" {
		t.Fatalf("unexpected json string list: %#v", got)
	}
}

func TestDataAgentSqlKeepsOnlyApprovedFields(t *testing.T) {
	raw, err := os.ReadFile("../../../storage/data/data.sql")
	if err != nil {
		t.Fatal(err)
	}
	sql := string(raw)

	for _, want := range []string{
		"CREATE TABLE IF NOT EXISTS `hg_data_agent`",
		"`agent_id` varchar(128) NOT NULL",
		"`register_status` varchar(16) NOT NULL DEFAULT 'pending'",
		"`dispatch_status` varchar(16) NOT NULL DEFAULT 'disabled'",
		"UNIQUE KEY `uk_data_agent_agent_id` (`agent_id`)",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("data.sql should contain %q", want)
		}
	}

	for _, forbidden := range []string{
		"certificate_fingerprint",
		"certificate_pem",
		"allowed_ip_cidrs",
		"capabilities",
		"total_slots",
		"available_slots",
		"online_status",
	} {
		if strings.Contains(sql, forbidden) {
			t.Fatalf("data.sql should not contain removed field %q", forbidden)
		}
	}
}
