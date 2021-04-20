package serializers

import(
	"github.com/aJuvan/mockAPI/serializers/text"
	"github.com/aJuvan/mockAPI/serializers/json"
)

// The base serializer struct.
// It contains a validation function, for the configuration; a serialization
// function, for serializing the data; and default content type, if it's
// not defined in the configuration.
type serializer struct {
	Validate func(data interface{}) bool
	Serialize func(data interface{}) string
	ContentType string
}

// List of all serializers
var Serializers map[string]serializer = map[string]serializer{
	"text": serializer{
		Validate: text.Validate,
		Serialize: text.Serialize,
		ContentType: text.ContentType,
	},
	"json": serializer{
		Validate: json.Validate,
		Serialize: json.Serialize,
		ContentType: json.ContentType,
	},
}
