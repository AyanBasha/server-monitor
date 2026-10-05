package main

import (
	"encoding/json"
	"fmt"

	//"io"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Name    string
	URL     string
	Timeout time.Duration
}

/**func getNetData(url string) int {
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
}**/

var (
	proxmox = Service{
		Name:    "Proxmox",
		URL:     "http://192.168.5.14:8096",
		Timeout: 5 * time.Second,
	}

	jellyfin = Service{
		Name:    "Jellyfin",
		URL:     "http://192.168.5.14:8096",
		Timeout: 5 * time.Second,
	}

	ha = Service{
		Name:    "Home Assistant",
		URL:     "http://192.168.5.14:8096",
		Timeout: 5 * time.Second,
	}

	crafty = Service{
		Name:    "Crafty",
		URL:     "http://192.168.5.14:8096",
		Timeout: 5 * time.Second,
	}

	immich = Service{
		Name:    "Immich",
		URL:     "http://192.168.5.14:8096",
		Timeout: 3 * time.Second,
	}

	piHole = Service{
		Name:    "Pi-hole",
		URL:     "http://192.168.5.5/admin",
		Timeout: 5 * time.Second,
	}
)

var services = []Service{proxmox, jellyfin, ha, crafty, immich, piHole}

func encode(services []Service) ([]byte, error) {
	b, err := json.MarshalIndent(services, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal services: %w", err)
	}

	return b, nil
}

func decode(b []byte) ([]Service, error) {
	var services []Service

	err := json.Unmarshal(b, &services)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal services %w", err)
	}

	return services, nil
}

func getUrlStatus(service Service) bool {
	client := http.Client{
		Timeout: service.Timeout,
	}

	req, err := http.NewRequest("HEAD", service.URL, nil)

	if err != nil {
		log.Fatal(err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func runServerMonitor() {

	fmt.Println("Server Monitor is starting...")

	for _, service := range services {

		fmt.Println("Service:", service.Name)
		fmt.Println("URL:", service.URL)
		fmt.Println("Up?:", getUrlStatus(service))
		//fmt.Println("Timeout:", service.Timeout)
		fmt.Println("\n")
	}
}

func main() {
	runServerMonitor()
}
