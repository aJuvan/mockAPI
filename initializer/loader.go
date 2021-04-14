package initializer

import(
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/aJuvan/mockAPI/serializers"
)

// Basic structs representing the configuration
type Router struct {
	Path string
	Method string
	Status int
	ContentType string `yaml:"content_type"`
	Serializer string
	Response map[interface{}]interface{}
}

type Config struct {
	Host string
	Prefix string
	Router []Router
}

// Load and validate the configuration
func LoadConfig(filename string) Config {

	// Open the configuration file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open configuration file.")
		os.Exit(2)
	}

	// Deserialize the contents
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot parse file.")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	// Normalize data
	config.Host = strings.TrimSpace(config.Host)
	config.Prefix = strings.TrimSpace(config.Prefix)
	for i, route := range config.Router {
		config.Router[i].Path = stringDefault(strings.TrimSpace(route.Path), "/")
		config.Router[i].Method = stringDefault(strings.ToUpper(strings.TrimSpace(route.Method)), "GET")
		config.Router[i].ContentType = strings.TrimSpace(route.ContentType)
		config.Router[i].Serializer = strings.TrimSpace(route.Serializer)
	}

	// Validate configuration or exit
	if success := validateConfig(config); !success {
		os.Exit(3)
	}

	return config
}

// Config validation marker
var configValid bool;

func validateConfig(config Config) bool {
	// Mark config as valid
	configValid = true

	// Test hostname
	configAssert(config.Host != "", "Host must not be empty.")

	// Test duplicated routes
	usedRoutes := make(map[string]bool)
	for i, route := range config.Router {

		// Create a route key as "METHOD PATH" and test it
		routePath := route.Method + " " + route.Path
		configAssert(!usedRoutes[routePath], "Route \"" + routePath + "\" already defined.")
		usedRoutes[routePath] = true

		// Check the status
		configAssert(route.Status != 0, "Return status not defined.")

		// Validate response, depending on the serializer
		if val, ok := serializers.Serializers[route.Serializer]; ok {
			val.Validate(route.Response)

			// Set the content type if empty
			if route.ContentType == "" {
				config.Router[i].ContentType = val.ContentType
			}
		} else {
			configAssert(false, "Must use one of the defined serializers.")
		}
	}

	return configValid
}

// Helper function for returning the default value where needed
func stringDefault(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// Helper function for printing errors on failed checks
func configAssert(test bool, errorMessage string) {
	if !test {
		fmt.Fprintln(os.Stderr, errorMessage)

		// Invalidate configuration if a test failed
		configValid = false
	}
}