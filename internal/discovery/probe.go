package discovery

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ManualProbe(ips []string) []Device {
	var devices []Device
	var mu sync.Mutex
	var wg sync.WaitGroup
	seen := make(map[string]bool)

	for _, ip := range ips {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			if d := ProbeTV(addr); d != nil {
				mu.Lock()
				if !seen[d.IP] {
					devices = append(devices, *d)
					seen[d.IP] = true
					fmt.Printf("Found: %s (%s)\n", d.Name, d.IP)
				}
				mu.Unlock()
			} else {
				fmt.Printf("No TV found at %s\n", addr)
			}
		}(ip)
	}

	wg.Wait()
	return devices
}

func ScanSubnet(myIP string) []Device {
	parts := strings.Split(myIP, ".")
	if len(parts) != 4 {
		return nil
	}
	subnet := strings.Join(parts[:3], ".")
	fmt.Printf("Scanning %s.0/24\n", subnet)
	
	var devices []Device
	var mu sync.Mutex
	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup
	seen := make(map[string]bool)
	
	for i := 1; i < 255; i++ {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			
			if d := ProbeTV(ip); d != nil {
				mu.Lock()
				if !seen[d.IP] {
					devices = append(devices, *d)
					seen[d.IP] = true
					fmt.Printf("Found: %s (%s)\n", d.Name, d.IP)
				}
				mu.Unlock()
			}
		}(fmt.Sprintf("%s.%d", subnet, i))
	}
	
	wg.Wait()
	return devices
}

func ProbeTV(ip string) *Device {
	ports := []string{"8008", "8009"}
	paths := []string{"/ssdp/device-desc.xml", "/dd.xml"}
	
	client := &http.Client{Timeout: 300 * time.Millisecond}
	
	for _, port := range ports {
		for _, path := range paths {
			url := fmt.Sprintf("http://%s:%s%s", ip, port, path)
			resp, err := client.Get(url)
			if err == nil && resp.StatusCode == 200 {
				resp.Body.Close()
				name := fetchDeviceName(url)
				return &Device{Name: name, IP: ip, Location: url}
			}
		}
	}
	return nil
}

func fetchDeviceName(location string) string {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(location)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()
	
	buf := make([]byte, 4096)
	n, _ := resp.Body.Read(buf)
	xmlStr := string(buf[:n])
	
	xmlStr = html.UnescapeString(xmlStr)
	
	start := strings.Index(xmlStr, "<friendlyName>")
	end := strings.Index(xmlStr, "</friendlyName>")
	if start != -1 && end != -1 {
		return xmlStr[start+14 : end]
	}
	return "Unknown"
}

func extractIP(url string) string {
	start := strings.Index(url, "://")
	if start == -1 {
		return ""
	}
	start += 3
	end := strings.Index(url[start:], ":")
	if end == -1 {
		end = strings.Index(url[start:], "/")
	}
	if end == -1 {
		return url[start:]
	}
	return url[start : start+end]
}
