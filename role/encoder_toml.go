package role

import "github.com/pelletier/go-toml"

// tomlMarshaler is a Marshaler that uses the go-toml package to marshal and unmarshal data.
type tomlMarshaler struct{}

// Marshal ...
func (tomlMarshaler) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}

// Unmarshal ...
func (tomlMarshaler) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}

// MarshalTOML ...
func (r *Roles) MarshalTOML() ([]byte, error) {
	return marshalRoles(r, tomlMarshaler{})
}

// UnmarshalTOML ...
func (r *Roles) UnmarshalTOML(b []byte) error {
	return unmarshalRoles(r, b, tomlMarshaler{})
}
