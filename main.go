package main

import(
	"fmt"
	"github.com/aJuvan/mockAPI/initializer"
	"github.com/aJuvan/mockAPI/server"
)

func main() {
	// Load settings
	settings := initializer.LoadSettings()

	// Load the configuration file
	config := initializer.LoadConfig(settings.Filename)

	// Return here if only a validation is needed
	fmt.Printf("%t\n", settings.Validate)
	if settings.Validate {
		return
	}

	// Serve the API
	server.Serve(settings, config)
}