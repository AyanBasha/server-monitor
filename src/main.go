package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Name    string
	URL     string
	Timeout time.Duration
}

func getNetData(url string) int {
	response, err := http.Get(url)

	if err != nil {
		log.Fatal(err)

		return -1
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		data := string(bodyBytes)
		fmt.Println(data)
	}

	return
}

var (
	proxmox = Service{
		Name:    "Proxmox",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	jellyfin = Service{
		Name:    "Jellyfin",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	ha = Service{
		Name:    "Home Assistant",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	crafty = Service{
		Name:    "Crafty",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	immich = Service{
		Name:    "Immich",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}
)

var services = []Service{proxmox, jellyfin, ha, crafty, immich}

func runServerMonitor() {

	fmt.Println("Server Monitor is starting...")

	for _, service := range services {

		fmt.Println("Service:", service.Name)
		fmt.Println("URL:", service.URL)
		fmt.Println("Timeout:", service.Timeout)
		fmt.Println("\n")
	}
}

func main() {
	runServerMonitor()
}
