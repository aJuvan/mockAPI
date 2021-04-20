package server

import(
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aJuvan/mockAPI/initializer"
	"github.com/aJuvan/mockAPI/serializers"
	"github.com/gorilla/mux"
)

// Prepare an http client for proxying requests
var httpClient = &http.Client{
  Timeout: time.Second * 10,
}

// Main serving function
func Serve(settings initializer.Settings, config initializer.Config) {
	// Define a new router
	router := mux.NewRouter()
	var subrouter *mux.Router
	
	// Create a subrouter for a host and optionally for the prefix
	if config.Prefix != "" {
		subrouter = router.
			Host(config.Host).
			PathPrefix(config.Prefix).
			Subrouter()
	} else {
		subrouter = router.
			Host(config.Host).
			Subrouter()
	}
	
	// Register each route and method
	for _, r := range config.Router {
		if (r.Proxy) {
			subrouter.
				HandleFunc(r.Path, proxyRoute(config.ProxyHost, r)).
				Methods(r.Method)
		} else {
			subrouter.
				HandleFunc(r.Path, serveRoute(r)).
				Methods(r.Method)
		}
	}

	// Serve the configured server
	fmt.Println("Listening on " + settings.Host + "...")
	http.ListenAndServe(settings.Host, router)
}

// Create a handler function for a configured route
func serveRoute(router initializer.Router) func (http.ResponseWriter, *http.Request) {
	// Get the response needed from the apropriate serializer
	serializer := serializers.Serializers[router.Serializer]
	response := serializer.Serialize(router.Response)

	// Create the handler function
	return func (w http.ResponseWriter, r *http.Request) {
		// Set the HTTP Status header
		w.WriteHeader(router.Status)

		// Set the Content type
		w.Header().Set("Content-Type", router.ContentType)

		// Return the calculated response
		fmt.Fprintln(w, response)
	}
}

// Create a handler function which proxies the request to the real API
func proxyRoute(proxyHost string, router initializer.Router) func (http.ResponseWriter, *http.Request) {
	// Prepare the server url
	proxyUrl := proxyHost + router.ProxyPath

	// Return the handler function
	return func (w http.ResponseWriter, r *http.Request) {
		// Update the url
		URL, _ := url.Parse(proxyUrl)
		r.URL.Scheme = URL.Scheme
		r.URL.Host = URL.Host
		r.URL.Path = proxyUrl
		r.RequestURI = ""
		r.Host = URL.Host

		// Set forwarding header
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)

		// Proxy the request
		if resp, err := httpClient.Do(r); err != nil {
			// In case of error, throw a 503
			w.WriteHeader(503)
			fmt.Fprintln(os.Stderr, "There was an error proxying the request:")
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(w, "Proxied service unavailable")
		} else {
			// Write each header
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			// Write status code
			w.WriteHeader(resp.StatusCode)

			// Copy response body
			io.Copy(w, resp.Body)
			resp.Body.Close()
		}
	}
}
