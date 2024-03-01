// Wraps the JSON library that we're using so that we
// can [easily??] change it

package json

import (
	stdlib "encoding/json"
	"io"

	json "github.com/goccy/go-json"
)

type Number = stdlib.Number
type RawMessage = stdlib.RawMessage

func Marshal(data any) ([]byte, error) {
	return json.Marshal(data)
}

func Pretty(data any) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func Unmarshal(data []byte, into any) error {
	return json.Unmarshal(data, into)
}

func MarshalInto(data any, w io.Writer) error {
	return json.NewEncoder(w).EncodeWithOption(data, json.UnorderedMap())
}

func NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}
