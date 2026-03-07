package discovery

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func SSDPSearch(searchTarget string, timeout time.Duration) []Device {
	remoteAddr, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	
	payload := fmt.Sprintf("M-SEARCH * HTTP/1.1\r\n"+
		"HOST:239.255.255.250:1900\r\n"+
		"MAN:\"ssdp:discover\"\r\n"+
		"MX:2\r\n"+
		"ST:%s\r\n"+
		"\r\n", searchTarget)
	
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil
	}
	defer conn.Close()
	
	conn.WriteTo([]byte(payload), remoteAddr)
	
	var devices []Device
	buf := make([]byte, 2048)
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}
		
		response := string(buf[:n])
		if loc, err := GetLocationURL(response); err == nil {
			ip := extractIP(loc)
			name := fetchDeviceName(loc)
			devices = append(devices, Device{Name: name, IP: ip, Location: loc})
		}
	}
	
	return devices
}

func GetLocationURL(response string) (string, error) {
	lines := strings.Split(response, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(line), "LOCATION:") {
			return strings.TrimSpace(line[9:]), nil
		}
	}
	return "", fmt.Errorf("LOCATION not found")
}
