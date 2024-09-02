package role

import "github.com/go-jose/go-jose/v3/json"

// jsonMarshaler is a Marshaler that uses the encoding/json package to marshal and unmarshal data.
type jsonMarshaler struct{}

// Marshal ...
func (jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (jsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MarshalJSON ...
func (r *Roles) MarshalJSON() ([]byte, error) {
	return marshalRoles(r, jsonMarshaler{})
}

// UnmarshalJSON ...
func (r *Roles) UnmarshalJSON(b []byte) error {
	return unmarshalRoles(r, b, jsonMarshaler{})
}

// MarshalBSON ...
func (r *Role) MarshalJSON() ([]byte, error) {
	return marshalSingularRole(r, jsonMarshaler{})
}

// UnmarshalBSON ...
func (r *Role) UnmarshalJSON(b []byte) error {
	return unmarshalSingularRole(r, b, jsonMarshaler{})
}
