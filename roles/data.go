package roles

import (
	"github.com/restartfu/gophig"
	"strings"
	"time"
)

// rolesData is a struct that is used to encode roles to BSON or any other format that requires encoding.
type rolesData struct {
	Roles       []string
	Expirations map[string]time.Time
}

func marshalRoles(r *Roles, marshaler gophig.Marshaler) ([]byte, error) {
	var d rolesData
	d.Expirations = make(map[string]time.Time)

	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	for _, rl := range r.roles {
		roleName := strings.ToLower(rl.Name())
		e, _ := r.roleExpirations[rl]
		if !e.IsZero() {
			d.Expirations[roleName] = e
		}
		d.Roles = append(d.Roles, roleName)
	}
	return marshaler.Marshal(d)
}

func unmarshalRoles(r *Roles, b []byte, marshaler gophig.Marshaler) error {
	var d rolesData
	if err := marshaler.Unmarshal(b, &d); err != nil {
		return err
	}

	rls := d.Roles
	for _, rl := range rls {
		ro, ok := ByName(rl)
		if ok {
			r.Add(ro)
			e, ok := d.Expirations[rl]
			if ok {
				r.Expire(ro, e)
			}
		}
	}
	return nil
}
