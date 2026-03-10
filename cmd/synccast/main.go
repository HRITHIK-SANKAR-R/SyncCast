package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

func main() {
	ipFlag := flag.String("ip", "", "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: synccast [options]\n\n  -ip string  Comma-separated TV IPs (e.g. --ip 192.168.1.42,192.168.1.50)\n")
	}
	flag.Parse()

	var manualIPs []string
	if *ipFlag != "" {
		for _, ip := range strings.Split(*ipFlag, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				manualIPs = append(manualIPs, ip)
			}
		}
	}

	fmt.Println("SyncCast: Discovering Android TVs...")

	devices := discovery.DiscoverTVs(manualIPs)

	if len(devices) == 0 {
		fmt.Println("No TVs found")
		return
	}

	fmt.Printf("\nFound %d device(s):\n", len(devices))
	for i, d := range devices {
		fmt.Printf("%d. %s (%s)\n", i+1, d.Name, d.IP)
	}
}
