package discovery_test

import (
	"net"
	"testing"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

func TestGetLocalIP_ReturnsValidIPv4(t *testing.T) {
	ip, err := discovery.GetLocalIP()
	if err != nil {
		t.Skipf("No local IP available (CI or isolated env): %v", err)
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		t.Fatalf("GetLocalIP returned invalid IP: %q", ip)
	}
	if parsed.To4() == nil {
		t.Fatalf("GetLocalIP returned non-IPv4 address: %q", ip)
	}
}

func TestGetLocalIP_NotLoopback(t *testing.T) {
	ip, err := discovery.GetLocalIP()
	if err != nil {
		t.Skipf("No local IP available: %v", err)
	}

	parsed := net.ParseIP(ip)
	if parsed.IsLoopback() {
		t.Fatalf("GetLocalIP returned loopback address: %q", ip)
	}
}
