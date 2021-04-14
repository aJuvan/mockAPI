package server

import(
	"fmt"
	"net/http"

	"github.com/aJuvan/mockAPI/initializer"
	"github.com/aJuvan/mockAPI/serializers"
	"github.com/gorilla/mux"
)

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
		subrouter.
			HandleFunc(r.Path, serveRoute(r)).
			Methods(r.Method)
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