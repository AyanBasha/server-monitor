package main

import (
	"fmt"
	"time"
)

type Service struct {
	Name    string
	URL     string
	Timeout time.Duration
}

func main() {
	jellyfin := Service{
		Name:    "Jellyfin",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	fmt.Println("Server Monitor is starting...")
	fmt.Println("Service:", jellyfin.Name)
	fmt.Println("URL:", jellyfin.URL)
	fmt.Println("Timeout:", jellyfin.Timeout)
}
