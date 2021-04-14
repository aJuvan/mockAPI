package text

import(
	"fmt"
	"os"
	"reflect"
)

// Default content type
var ContentType string = "text/html"

// Configuration response validation function
func Validate(data map[interface{}]interface{}) bool {
	// Check if the "string" key is present in the response
	if val, ok := data["string"]; !ok {
		fmt.Fprintln(os.Stderr, "The parameter \"string\" is missing in the response")
		return false
	// Check if the "string" key is of type string
	} else if reflect.TypeOf(val).Kind() != reflect.String {
		fmt.Fprintln(os.Stderr, "The parameter \"string\" must be of type string")
		return false
	}

	return true
}

// Return the defined string
func Serialize(data map[interface{}]interface{}) string {
	return data["string"].(string);
}