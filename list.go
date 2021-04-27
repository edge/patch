package patch

import (
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"
)

// List is a map of patches keyed by version number.
// This makes it simple to coordinate any number of patches and ensure they run in a predictable order.
type List map[string]Patch

// Apply patches.
func (pl List) Apply(cs version.Constraints) error {
	vers, err := pl.Versions()
	if err != nil {
		return err
	}
	for _, ver := range vers {
		if cs.Check(ver) {
			patch, ok := pl[ver.String()]
			if !ok {
				return fmt.Errorf("unable to locate patch \"%s\"", ver.String())
			}
			if err := patch.Apply(); err != nil {
				return &Error{InVersion: ver, IsApply: true, OriginalError: err}
			}
		}
	}
	return nil
}

// HighestVersion defined in the patch list.
func (pl List) HighestVersion() (*version.Version, error) {
	vers, err := pl.Versions()
	if err != nil {
		return nil, err
	}
	return vers[len(vers)-1], nil
}

// Pick versions.
func (pl List) Pick(cs version.Constraints) (version.Collection, error) {
	vers, err := pl.Versions()
	if err != nil {
		return vers, err
	}
	picked := version.Collection{}
	for _, ver := range vers {
		if cs.Check(ver) {
			picked = append(picked, ver)
		}
	}
	if len(picked) == 0 {
		return picked, fmt.Errorf("no patches for %s", cs.String())
	}
	return picked, nil
}

// Revert patches.
func (pl List) Revert(cs version.Constraints) error {
	vers, err := pl.Pick(cs)
	if err != nil {
		return err
	}
	i := len(vers)
	for {
		i--
		ver := vers[i]
		if cs.Check(ver) {
			patch, ok := pl[ver.String()]
			if !ok {
				return fmt.Errorf("unable to locate patch \"%s\"", ver.String())
			}
			if err := patch.Revert(); err != nil {
				return &Error{InVersion: ver, IsApply: false, OriginalError: err}
			}
		}
		if i == 0 {
			break
		}
	}
	return nil
}

// Versions in the patch list.
// Produces a collection (slice) of versions from lowest to highest.
// An error is returned if the patch list is empty or a version is invalid.
func (pl List) Versions() ([]*version.Version, error) {
	if len(pl) == 0 {
		return nil, errors.New("patch list empty")
	}
	vers := make([]*version.Version, len(pl))
	i := 0
	for vstr := range pl {
		ver, err := version.NewVersion(vstr)
		if err != nil {
			return vers, err
		}
		vers[i] = ver
		i++
	}
	sort.Sort(version.Collection(vers))
	return vers, nil
}
