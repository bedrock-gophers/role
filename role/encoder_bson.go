package role

import (
	"github.com/rcrowley/go-bson"
)

// bsonMarshaler is a Marshaler that uses the go-bson package to marshal and unmarshal data.
type bsonMarshaler struct{}

// Marshal ...
func (bsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return bson.Marshal(v)
}

// Unmarshal ...
func (bsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return bson.Unmarshal(data, v)
}

// MarshalBSON ...
func (r *Roles) MarshalBSON() ([]byte, error) {
	return marshalRoles(r, bsonMarshaler{})
}

// UnmarshalBSON ...
func (r *Roles) UnmarshalBSON(b []byte) error {
	return unmarshalRoles(r, b, bsonMarshaler{})
}
