package text

import(
	"fmt"
	"os"
	"reflect"
)

// Default content type
var ContentType string = "text/html"

// Configuration response validation function
func Validate(data interface{}) bool {
	// Check if the data is of type string
	if reflect.TypeOf(data).Kind() != reflect.String {
		fmt.Fprintln(os.Stderr, "The response must be of type string")
		return false
	}

	return true
}

// Return the defined string
func Serialize(data interface{}) string {
	return data.(string);
}
