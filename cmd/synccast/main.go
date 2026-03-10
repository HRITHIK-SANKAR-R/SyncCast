package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/streamer"
)

func main() {
	ipFlag := flag.String("ip", "", "")
	fileFlag := flag.String("file", "", "")
	portFlag := flag.Int("port", 6969, "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: synccast [options]\n\n"+
			"  -ip string    Comma-separated TV IPs (e.g. --ip 192.168.1.42,192.168.1.50)\n"+
			"  -file string  Path to media file to stream\n"+
			"  -port int     HTTP server port (default 6969)\n")
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

	if *fileFlag == "" {
		return
	}

	srv, err := streamer.New(*fileFlag, *portFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	myIP, _ := discovery.GetLocalIP()
	fmt.Printf("\nStream URL: %s\n", srv.StreamURL(myIP))

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
