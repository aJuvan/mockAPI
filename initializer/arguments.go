package initializer

import(
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Structure holding all of the command line parsed settings.
// Fetch them by using loadSettings()
type Settings struct {
	Validate bool
	Host string
	Filename string
}

// Load program settings and return a populated Settings struct
func LoadSettings() Settings {
	settings := Settings{}

	// Define available settings and their defaults
	validate := flag.Bool("validate", false, "Only validate configuration.")
	host := flag.String("host", getEnvOrDefault("MOCKAPI__SETTINGS_HOST", "127.0.0.1:8080"), "Hostname listen address.")

	// Parse the flags
	flag.Parse()

	// Save the parsed data
	settings.Validate = *validate
	settings.Host = *host

	// Asure the positional of the arguments are present
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Expected only a single file argument.")
		os.Exit(1)
	}

	// Load the rest of the arguments
	settings.Filename = flag.Args()[0]

	return settings
}

func getEnvOrDefault(key string, defaultValue string) string {
	// Get environment value
	value := os.Getenv(key)

	// If the value is empty, return the default
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	// Get environment value
	value := os.Getenv(key)

	// If the value is empty return the default
	if len(value) == 0 {
		return defaultValue
	}

	// Convert string to int and exit if there was an error
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Environment variable %s is not an integer.", key))
		os.Exit(1)
	}

	return intValue
}