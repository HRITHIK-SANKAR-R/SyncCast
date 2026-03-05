package main

import (
	"fmt"
	"log"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

func main() {
	fmt.Println("Synccast : Initialising")
	
	ip, err := discovery.GetLocalIP(); 
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current ip of the wifi: %s\n", ip)
	
}
