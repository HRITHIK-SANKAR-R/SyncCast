package discovery

import (
	"fmt"
	"sync"
	"time"
)

func DiscoverAllTVs() []Device {
	myIP, _ := GetLocalIP()
	fmt.Printf("Your IP: %s\n", myIP)
	
	var devices []Device
	var mu sync.Mutex
	var wg sync.WaitGroup
	seen := make(map[string]bool)

	fmt.Println("Trying SSDP discovery...")
	targets := []string{
		"urn:dial-multiscreen-org:service:dial:1",
		"urn:schemas-upnp-org:device:MediaRenderer:1",
	}

	for _, target := range targets {
		wg.Add(1)
		go func(st string) {
			defer wg.Done()
			found := SSDPSearch(st, 2*time.Second)
			mu.Lock()
			for _, d := range found {
				if !seen[d.IP] {
					devices = append(devices, d)
					seen[d.IP] = true
				}
			}
			mu.Unlock()
		}(target)
	}

	wg.Wait()
	
	if len(devices) == 0 {
		fmt.Println("\nSSDP failed, scanning subnet...")
		devices = ScanSubnet(myIP)
	}
	
	return devices
}
