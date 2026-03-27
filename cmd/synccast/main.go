package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/control"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/state"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/streamer"
)

func main() {
	ipFlag := flag.String("ip", "", "Comma-separated TV IPs (e.g. --ip 192.168.1.42)")
	fileFlag := flag.String("file", "", "Path to media file to stream (optional — can be set from dashboard)")
	portFlag := flag.Int("port", 6969, "HTTP server port (default 6969)")
	flag.Parse()

	myIP, _ := discovery.GetLocalIP()

	// Build server — file is optional now
	srv := streamer.NewServer(*portFlag)
	srv.HostIP = myIP

	// If a file was provided via CLI, set it immediately
	if *fileFlag != "" {
		if err := srv.SetFile(*fileFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}

	// Run device discovery in the background (non-blocking)
	var manualIPs []string
	if *ipFlag != "" {
		for _, ip := range strings.Split(*ipFlag, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				manualIPs = append(manualIPs, ip)
			}
		}
	}

	go func() {
		fmt.Println("SyncCast: Discovering devices...")
		devices := discovery.DiscoverTVs(manualIPs)
		srv.SetDevices(devices)
		if len(devices) > 0 {
			fmt.Printf("Found %d device(s):\n", len(devices))
			for i, d := range devices {
				fmt.Printf("  %d. %s (%s)\n", i+1, d.Name, d.IP)
			}
		} else {
			fmt.Println("No devices found — scan from dashboard or use --ip")
		}
	}()

	// Wire up WebSocket hub and state manager
	hub := control.NewHub()
	stateManager := state.NewManager()
	srv.StateSnapshot = func() interface{} {
		return stateManager.Snapshot()
	}
	hub.SetLifecycleHooks(
		func(role control.Role) {
			stateManager.OnClientConnected(string(role))
			s := stateManager.Snapshot()
			log.Printf("state: %s (remote=%d player=%d)", s.State, s.RemoteClients, s.PlayerClients)
		},
		func(role control.Role) {
			stateManager.OnClientDisconnected(string(role))
			s := stateManager.Snapshot()
			log.Printf("state: %s (remote=%d player=%d)", s.State, s.RemoteClients, s.PlayerClients)
		},
	)
	go hub.Run()
	srv.WSHandler = hub.HandleWS

	// Attach discovery function so dashboard can trigger scans
	srv.DiscoverFunc = func(ips []string) []discovery.Device {
		return discovery.DiscoverTVs(ips)
	}

	fmt.Printf("\n✓ SyncCast running on port %d\n", *portFlag)
	fmt.Printf("  Dashboard:  http://%s:%d/dashboard\n", myIP, *portFlag)
	fmt.Printf("  Remote:     http://%s:%d/remote\n", myIP, *portFlag)
	fmt.Printf("  Player:     http://%s:%d/player\n", myIP, *portFlag)
	fmt.Println()

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
