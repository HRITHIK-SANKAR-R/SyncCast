package main

import (
	"fmt"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

func main() {
	fmt.Println("SyncCast: Discovering Android TVs...")
	
	devices := discovery.DiscoverAllTVs()
	
	if len(devices) == 0 {
		fmt.Println("No TVs found")
		return
	}
	
	fmt.Printf("\nFound %d device(s):\n", len(devices))
	for i, d := range devices {
		fmt.Printf("%d. %s (%s)\n", i+1, d.Name, d.IP)
	}
}
