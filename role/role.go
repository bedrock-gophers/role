package role

import (
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Role represents a role in the game.
type Role struct {
	// name is the name of the role.
	name string
	// inherits is a role that this role inherits from.
	inherits string
	// colour is the colour of the role.
	colour string
	// tier is the tier of the role.
	tier int
}

// Name returns the name of the role.
func (r Role) Name() string {
	return r.name
}

// Inherits returns the role that this role inherits from.
func (r Role) Inherits() (Role, bool) {
	rl, ok := rolesName[r.inherits]
	return rl, ok
}

// Tier returns the tier of the role.
func (r Role) Tier() int {
	return r.tier
}

// Coloured returns the given string coloured with the role's colour.
func (r Role) Coloured(s string) string {
	if len(r.colour) == 0 {
		return s
	}
	if strings.HasPrefix(r.colour, "§") {
		return r.colour + s + "§r"
	}
	return text.Colourf("<%s>%s</%s>", r.colour, s, r.colour)
}

// Roles manages a list of role for a user. Roles can be added, removed, and checked for. Roles can also have an
// expiration time, after which they are removed from the user's role list.
type Roles struct {
	roleMu          sync.Mutex
	roles           []Role
	roleExpirations map[Role]time.Time
}

// NewRoles creates a new Roles instance.
func NewRoles(roles []Role, expirations map[Role]time.Time) *Roles {
	r := &Roles{
		roles:           roles,
		roleExpirations: expirations,
	}
	r.removeDuplicates()
	return r
}

// Add adds a role to the manager's role list.
func (r *Roles) Add(ro Role) {
	if r.Contains(ro) {
		return
	}
	r.roleMu.Lock()
	r.roles = append(r.roles, ro)
	r.roleMu.Unlock()
	sortRoles(r.roles)
}

// Remove removes a role from the manager's role list. Users are responsible for updating the highest role usages if
// changed.
func (r *Roles) Remove(ro Role) bool {
	r.roleMu.Lock()
	i := slices.IndexFunc(r.roles, func(other Role) bool {
		return ro == other
	})
	r.roles = slices.Delete(r.roles, i, i+1)
	delete(r.roleExpirations, ro)
	r.roleMu.Unlock()
	r.checkExpiry()
	sortRoles(r.roles)
	return true
}

// Contains returns true if the manager has any of the given role. Users are responsible for updating the highest role
// usages if changed.
func (r *Roles) Contains(roles ...Role) bool {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	var actualRoles []Role
	for _, ro := range r.roles {
		r.propagateRoles(&actualRoles, ro)
	}

	for _, r := range roles {
		if i := slices.IndexFunc(actualRoles, func(other Role) bool {
			return r == other
		}); i >= 0 {
			return true
		}
	}
	return false
}

// Expiration returns the expiration time for a role. If the role does not expire, the second return value will be false.
func (r *Roles) Expiration(ro Role) (time.Time, bool) {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	e, ok := r.roleExpirations[ro]
	return e, ok
}

// Expire sets the expiration time for a role. If the role does not expire, the second return value will be false.
func (r *Roles) Expire(ro Role, t time.Time) {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	r.roleExpirations[ro] = t
}

// Highest returns the highest role the manager has, in terms of hierarchy.
func (r *Roles) Highest() Role {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	return r.roles[len(r.roles)-1]
}

// All returns the user's role.
func (r *Roles) All() []Role {
	r.checkExpiry()
	r.roleMu.Lock()
	defer r.roleMu.Unlock()
	return append(make([]Role, 0, len(r.roles)), r.roles...)
}

// removeDuplicates removes duplicate roles from the user's role list.
func (r *Roles) removeDuplicates() {
	r.roleMu.Lock()
	defer r.roleMu.Unlock()

	var rls []Role
	seen := map[Role]struct{}{}
	for _, ro := range r.roles {
		if _, ok := seen[ro]; !ok {
			seen[ro] = struct{}{}
			rls = append(rls, ro)
		}
	}
	r.roles = rls

}

// propagateRoles propagates roles to the user's role list.
func (r *Roles) propagateRoles(actualRoles *[]Role, role Role) {
	*actualRoles = append(*actualRoles, role)
	if h, ok := role.Inherits(); ok {
		r.propagateRoles(actualRoles, h)
	}
}

// checkExpirations checks each role the user has and removes the expired ones.
func (r *Roles) checkExpiry() {
	r.roleMu.Lock()
	rl, expirations := r.roles, r.roleExpirations
	r.roleMu.Unlock()

	for _, ro := range rl {
		if t, ok := expirations[ro]; ok && time.Now().After(t) {
			r.Remove(ro)
		}
	}
}
