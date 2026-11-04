package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"crypto/tls"
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

type Service struct {
	Name          string        `json:"name"`
	URL           string        `json:"url"`
	Timeout       time.Duration `json:"timeout"`
	SkipTLSVerify bool          `json:"skip_tls_verify"`
}

// services is shared between HTTP handlers. Go's http server runs every
// request in its own goroutine, so two requests can touch this slice at the
// same time. servicesMu makes sure only one goroutine reads/writes at once.
var (
	servicesMu sync.Mutex
	services   = []Service{
		{Name: "Proxmox", URL: "https://192.168.5.26:8006", Timeout: 5 * time.Second, SkipTLSVerify: true},
		{Name: "Jellyfin", URL: "http://192.168.5.14:8096", Timeout: 5 * time.Second},
		{Name: "Home Assistant", URL: "http://192.168.5.9:80", Timeout: 5 * time.Second},
		{Name: "Pi-hole", URL: "http://192.168.5.5/admin", Timeout: 5 * time.Second},
	}
)

// addService appends under the lock.
func addService(s Service) {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	services = append(services, s)
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

// ---------------------------------------------------------------------------
// Checking
// ---------------------------------------------------------------------------

func getUrlStatus(service Service) bool {
	client := http.Client{Timeout: service.Timeout}



	

	req, err := http.NewRequest(http.MethodHead, service.URL, nil)
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

	log.Printf("[%s] HTTP %d", service.Name, resp.StatusCode)

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func runServerMonitor(w io.Writer) {
	fmt.Fprintln(w, "=+=+= Server Monitor =+=+=")
	fmt.Fprintln(w)

	for _, service := range snapshotServices() {
		fmt.Fprintln(w, "Service:", service.Name)
		fmt.Fprintln(w, "URL:", service.URL)
		fmt.Fprintln(w, "Accessible?:", getUrlStatus(service))
		fmt.Fprintln(w)
	}
}

func dashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	runServerMonitor(w)
}

// GET /form -> show the HTML form
func formPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, templatesDir+"/form.html")
}

// POST /form -> validate, add the service, redirect to the dashboard
func formSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// These keys must match the name="" attributes in form.html.
	name := strings.TrimSpace(r.FormValue("name"))
	rawURL := strings.TrimSpace(r.FormValue("url"))

	if name == "" || rawURL == "" {
		http.Error(w, "name and url are both required", http.StatusBadRequest)
		return
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "invalid url: "+rawURL, http.StatusBadRequest)
		return
	}

	addService(Service{Name: name, URL: u.String(), Timeout: 5 * time.Second})
	log.Printf("Added service %q -> %s", name, u.String())

	// 303 See Other tells the browser: "now GET this other page".
	// This is the Post/Redirect/Get pattern — refreshing the dashboard
	// won't re-submit the form.
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func main() {
	mux := http.NewServeMux()

	// Go 1.22+ patterns: "METHOD /path". The router rejects other methods
	// with 405 for us, so the handlers don't need manual method checks.
	mux.HandleFunc("GET /dashboard", dashHandler)
	mux.HandleFunc("GET /form", formPageHandler)
	mux.HandleFunc("POST /form", formSubmitHandler)

	// "/" is the catch-all: anything not matched above is looked up as a
	// file in ./templates. "/" -> templates/index.html automatically.
	mux.Handle("/", http.FileServer(http.Dir(templatesDir)))

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
