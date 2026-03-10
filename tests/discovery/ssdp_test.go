package discovery_test

import (
	"testing"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

func TestGetLocationURL_ValidResponse(t *testing.T) {
	response := "HTTP/1.1 200 OK\r\n" +
		"CACHE-CONTROL: max-age=1800\r\n" +
		"LOCATION: http://192.168.1.100:8008/ssdp/device-desc.xml\r\n" +
		"ST: urn:dial-multiscreen-org:service:dial:1\r\n" +
		"\r\n"

	loc, err := discovery.GetLocationURL(response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := "http://192.168.1.100:8008/ssdp/device-desc.xml"
	if loc != expected {
		t.Fatalf("expected %q, got %q", expected, loc)
	}
}

func TestGetLocationURL_CaseInsensitive(t *testing.T) {
	response := "HTTP/1.1 200 OK\r\n" +
		"location: http://10.0.0.5:8009/dd.xml\r\n" +
		"\r\n"

	loc, err := discovery.GetLocationURL(response)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := "http://10.0.0.5:8009/dd.xml"
	if loc != expected {
		t.Fatalf("expected %q, got %q", expected, loc)
	}
}

func TestGetLocationURL_Missing(t *testing.T) {
	response := "HTTP/1.1 200 OK\r\n" +
		"ST: urn:dial-multiscreen-org:service:dial:1\r\n" +
		"\r\n"

	_, err := discovery.GetLocationURL(response)
	if err == nil {
		t.Fatal("expected error for missing LOCATION header")
	}
}

func TestGetLocationURL_EmptyResponse(t *testing.T) {
	_, err := discovery.GetLocationURL("")
	if err == nil {
		t.Fatal("expected error for empty response")
	}
}
