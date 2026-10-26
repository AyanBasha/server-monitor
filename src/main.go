package main

import (
	"encoding/json"
	"fmt"
	"net"

	//"io"
	"log"
	"net/http"
	"os"
	"time"
	//"://tailscale.com"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface {}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		return err
	default:
		return fmt.Errorf("invalid duration type")
	}
}

type Service struct {
	Name    string         `json:"name"`
	URL     string         `json:"url"`
	Timeout time.Duration  `json:"timeout"`
}

/**func loadServices(filename String) ([]String, error) {
	
}**/

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
		log.Printf("Error creating request %s: %v", service.Name, err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error checking status %s: %v", service.Name, err)
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
		fmt.Println("Accessible?:", getUrlStatus(service))
		//fmt.Println("Timeout:", service.Timeout)
		fmt.Println("\n")
	}
}

func dashHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET"{
		http.Error(w, "method is not supported", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "You have reached your dashboard")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("Service Name")
	address := r.FormValue("IP Address")
	fmt.Fprintf(w, "Service Name = %s\n", name)
	fmt.Fprintf(w, "IP Address = %s\n", address)

}

func main() {
	//var httpClient *http.Client
	var ln net.Listener
	var err error

	runServerMonitor()

	

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/dashboard", dashHandler)

	fmt.Println("Starting server at port 8080")
	
	if err := http.ListenAndServe(":8080", nil); err !=nil {
		log.Fatal(err)
	}
	
	if os.Getenv("TAILSCALE_ENABLED") == "true" {
		fmt.Println("Initializing embedded Tailscale network...")

		srv := &tsnet.Server{
			Hostname: "universal-checker",
			AuthKey: os.Getenv("TS_AUTHKEY"),
		}

		defer srv.close()

		httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
					return srv.Dial(ctx, network, address)
				},
			},
			Timeout: 5*time.Second,
		}
	} else {
		fmt.Println("Using standard system network...")
		
		httpClient = &http.Client{
			Timeout: 5*time.Second,
		}
	}

}
