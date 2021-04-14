package serializers

import(
	"github.com/aJuvan/mockAPI/serializers/text"
)

// The base serializer struct.
// It contains a validation function, for the configuration; a serialization
// function, for serializing the data; and default content type, if it's
// not defined in the configuration.
type serializer struct {
	Validate func(data map[interface{}]interface{}) bool
	Serialize func(data map[interface{}]interface{}) string
	ContentType string
}

// List of all serializers
var Serializers map[string]serializer = map[string]serializer{
	"text": serializer{
		Validate: text.Validate,
		Serialize: text.Serialize,
		ContentType: text.ContentType,
	},
}