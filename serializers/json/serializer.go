package json

import(
	"fmt"
	"encoding/json"
	"os"
)

// Default content type
var ContentType string = "application/json"

// Configuration response validation function
func Validate(data interface{}) bool {
	// Check if the data can be serialized
	if _, err := json.Marshal(data); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot serialize response as json")
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	return true
}

// Return the serialized data
func Serialize(data interface{}) string {
	val, _ := json.Marshal(data)
	return string(val);
}
