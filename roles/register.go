package roles

import (
	"errors"
	"fmt"
	"github.com/restartfu/gophig"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	// roleMu is a mutex that protects the roles slice.
	roleMu sync.Mutex
	// roles is a slice of all roles.
	roles []Role
	// rolesName is a map of all roles.
	rolesName = map[string]Role{}
)

// register registers a role.
func register(rls ...Role) {
	roleMu.Lock()
	for _, r := range rls {
		rolesName[strings.ToLower(r.Name())] = r
		roles = append(roles, r)
	}
	roleMu.Unlock()
}

// Load loads all roles from a folder.
func Load(folder string) error {
	folder = strings.TrimSuffix(folder, "/")
	files, err := os.ReadDir(folder)
	if err != nil {
		return errors.New(fmt.Sprintf("error loading roles: %v", err))
	}

	var newRoles []Role
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		r, err := loadRole(folder + "/" + file.Name())
		if err != nil {
			return errors.New(fmt.Sprintf("error loading role %s: %v", file.Name(), err))
		}
		newRoles = append(newRoles, r)
	}

	roleMu.Lock()
	roles = make([]Role, 0)
	rolesName = map[string]Role{}
	roleMu.Unlock()

	sortRoles(newRoles)
	register(newRoles...)
	return nil
}

// roleData is a struct that is used to decode roles from JSON.
type roleData struct {
	Name     string `json:"name"`
	Inherits string `json:"inherits,omitempty"`
	Colour   string `json:"colour,omitempty"`
	Tier     int    `json:"tier"`
}

// loadRole loads a role from a file.
func loadRole(filePath string) (Role, error) {
	var data roleData
	err := gophig.GetConfComplex(filePath, gophig.JSONMarshaler{}, &data)
	if err != nil {
		return Role{}, err
	}

	for _, r := range roles {
		if r.Name() == data.Name {
			return Role{}, errors.New("role with name " + data.Name + " already exists")
		}
		if r.tier == data.Tier {
			return Role{}, errors.New(fmt.Sprintf("role with tier %d already exists", data.Tier))
		}
	}

	return Role{
		name:     data.Name,
		inherits: data.Inherits,
		colour:   data.Colour,
		tier:     data.Tier,
	}, nil
}

// sortRoles sorts the roles by their tier.
func sortRoles(rls []Role) {
	roleMu.Lock()
	sort.SliceStable(rls, func(i, j int) bool {
		return rls[i].tier < rls[j].tier
	})
	roleMu.Unlock()
}

// All returns all roles that are currently registered.
func All() []Role {
	roleMu.Lock()
	r := make([]Role, len(roles))
	copy(r, roles)
	roleMu.Unlock()
	return r
}

// ByName returns a role by its name.
func ByName(name string) (Role, bool) {
	roleMu.Lock()
	r, ok := rolesName[strings.ToLower(name)]
	roleMu.Unlock()
	return r, ok
}
