package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

// listenAddr is where the web server listens. ":8080" means "every network
// interface on port 8080", so other machines on your LAN can reach the LXC.
const listenAddr = ":8090"

// templatesDir is RELATIVE TO THE DIRECTORY YOU RUN THE PROGRAM FROM.
// air runs from the project root, so "./templates" is right.
// If you ever run `go run .` from inside src/, this would need "../templates".
const templatesDir = "./templates"

// ---------------------------------------------------------------------------
// Data model
// ---------------------------------------------------------------------------

type status int

const (
	down status = iota
	unusual
	up
	checkConfig
)

func (s status) String() string {
	switch s {
	case down:
		return "Service is down"
	case unusual:
		return "Up, but unusual"
	case up:
		return "Service is up"
	case checkConfig:
		return "check service configuration"
	}

	return "unknown status"
}

func (s status) Class() string {
	switch s {
	case up:
		return "up"
	case unusual:
		return "unusual"
	case checkConfig:
		return "config"
	}
	return "down"
}

type Service struct {
	Name          string        `json:"name"`
	URL           string        `json:"url"`
	Timeout       time.Duration `json:"timeout"`
	SkipTLSVerify bool          `json:"skip_tls_verify"`
	ID            int           `json:"ID"`
}

type CheckResult struct {
	Service Service
	Status  status
	Message string
}

// services is shared between HTTP handlers. Go's http server runs every
// request in its own goroutine, so two requests can touch this slice at the
// same time. servicesMu makes sure only one goroutine reads/writes at once.
var (
	servicesMu sync.Mutex
	services   []Service
	nextID     int
)

// addService appends under the lock.
func addService(s Service) {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	s.ID = nextID
	nextID++
	services = append(services, s)
	saveLocked()
}

// snapshotServices returns a COPY of the slice. We hold the lock only long
// enough to copy, then release it — so slow network checks (up to 5s each)
// never block someone submitting the form.
func snapshotServices() []Service {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	out := make([]Service, len(services))
	copy(out, services)
	return out
}

func checkAll() []CheckResult {
	var results []CheckResult

	for _, service := range snapshotServices() {
		st, msg := getUrlStatus(service)
		results = append(results, CheckResult{
			Service: service,
			Status:  st,
			Message: msg,
		})
	}

	return results
}

func deleteService(id int) bool {
	servicesMu.Lock()
	defer servicesMu.Unlock()

	for i, s := range services {
		if s.ID == id {
			services = append(services[:i], services[i+1:]...)
			saveLocked()
			return true
		}
	}
	return false
}

//HERER

func serviceFromForm(r *http.Request) (Service, error) {
	if err := r.ParseForm(); err != nil {
		return Service{}, fmt.Errorf("could not parse form: %w", err)
	}

	name := strings.TrimSpace(r.FormValue("name"))
	rawURL := strings.TrimSpace(r.FormValue("url"))

	if name == "" || rawURL == "" {
		return Service{}, fmt.Errorf("name and url are both required")
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Service{}, fmt.Errorf("invalid url: %s", rawURL)
	}

	skipTLS := "" != r.FormValue("skipTLS")

	return Service{
		Name:          name,
		URL:           u.String(),
		Timeout:       5 * time.Second,
		SkipTLSVerify: skipTLS,
	}, nil
}

func updateService(id int, updated Service) bool {
	servicesMu.Lock()
	defer servicesMu.Unlock()

	for i, s := range services {
		if s.ID == id {
			updated.ID = id
			services[i] = updated
			saveLocked()
			return true
		}
	}

	return false
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Which service? (copy the first part of deleteHandler)
    id, err := strconv.Atoi(r.FormValue("ID"))
    if err != nil {
        http.Error(w, "invalid service ID", http.StatusBadRequest)
        return
    }

    // 2. What are the new values? (like formSubmitHandler)
    updated, err := serviceFromForm(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 3. Apply the change
    if !updateService(id, updated) {
        http.Error(w, "service not found", http.StatusNotFound)
        return
    }

    log.Printf("Updated service %d -> %q %s", id, updated.Name, updated.URL)
    http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// ---------------------------------------------------------------------------
// Checking
// ---------------------------------------------------------------------------

func getUrlStatus(service Service) (status, string) {
	client := http.Client{Timeout: service.Timeout}

	if service.SkipTLSVerify {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	req, err := http.NewRequest(http.MethodGet, service.URL, nil)
	if err != nil {
		log.Printf("Error creating request %s: %v", service.Name, err)
		return checkConfig, fmt.Sprintf("Error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error checking status %s: %v", service.Name, err)
		return down, fmt.Sprintf("Error checking status: %v", err)
	}
	defer resp.Body.Close()
	//put a dropdown next to the status up or down and output this message

	result := down

	if resp.StatusCode < 400 {
		result = up
	} else if resp.StatusCode < 500 {
		result = unusual
	} else if resp.StatusCode >= 500 {
		result = down
	}

	return result, fmt.Sprintf("HTTP %d", resp.StatusCode)
}

func runServerMonitor(w io.Writer) {
	fmt.Fprintln(w, "=+=+= Server Monitor =+=+=")
	fmt.Fprintln(w)

	for _, service := range snapshotServices() {
		fmt.Fprintln(w, "Service:", service.Name)
		fmt.Fprintln(w, "URL:", service.URL)

		serviceStatus, detail := getUrlStatus(service)
		fmt.Fprintln(w, "Status:", serviceStatus)
		fmt.Fprintln(w, detail)

		fmt.Fprintln(w)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "invalid service ID", http.StatusBadRequest)
		return
	}

	if !deleteService(id) {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func dashHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Load the template
	tmpl, err := template.ParseFiles(templatesDir + "/dashboard.html")
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Run the checks
	results := checkAll()

	// 3. Tell the browser it's HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 4. Fill in the template and send it
	if err := tmpl.Execute(w, results); err != nil {
		log.Printf("template execute error: %v", err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func formPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, templatesDir+"/form.html")
}

func formSubmitHandler(w http.ResponseWriter, r *http.Request) {
	s, err := serviceFromForm(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addService(s)
	log.Printf("Added service %q -> %s", s.Name, s.URL)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func main() {
	if err := loadServices(); err != nil {
		log.Fatalf("could not load %s: %v", dataFile, err)
	}

	mux := http.NewServeMux()

	// Go 1.22+ patterns: "METHOD /path". The router rejects other methods
	// with 405 for us, so the handlers don't need manual method checks.
	mux.HandleFunc("GET /{$}", rootHandler)
	mux.HandleFunc("GET /dashboard", dashHandler)
	mux.HandleFunc("GET /form", formPageHandler)
	mux.HandleFunc("POST /form", formSubmitHandler)
	mux.HandleFunc("POST /delete", deleteHandler)
	mux.HandleFunc("POST /edit", editHandler)

	   mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Run the terminal check in the BACKGROUND (a goroutine). Before, this ran
	// first and could block startup for 5s per unreachable service — painful
	// with air restarting on every save.
	go func() {
		fmt.Println("Running check in terminal")
		runServerMonitor(os.Stdout)
	}()

	fmt.Println("Starting server on", listenAddr)

	// ListenAndServe blocks forever (until error). Any code AFTER this line
	// never runs — which is why the old Tailscale block was dead code.
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

// TODO(tailscale): the tsnet setup that was here is parked. When you come
// back to it, it has to run BEFORE ListenAndServe, needs imports for
// "context" and "tailscale.com/tsnet" (`go get tailscale.com/tsnet`), and the
// resulting httpClient must actually be used by getUrlStatus.
