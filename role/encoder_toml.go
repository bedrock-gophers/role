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

// MarshalBSON ...
func (r *Role) MarshalTOML() ([]byte, error) {
	return marshalSingularRole(r, tomlMarshaler{})
}

// UnmarshalBSON ...
func (r *Role) UnmarshalTOML(b []byte) error {
	return unmarshalSingularRole(r, b, tomlMarshaler{})
}
