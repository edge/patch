package patch

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

type TestPatch struct {
	E error
}

func (p TestPatch) Apply() error {
	return p.E
}

func (p TestPatch) Revert() error {
	return p.E
}

var errNoReason = errors.New("no reason")

var testOrd []string = []string{"0.1.0", "1.0.0", "1.7.3", "2.0.1", "64.3.2"}

// versions are jumbled in this test to ensure failure case if patch list doesn't sort correctly.
// Please don't do this in actual code!
var testPatchList List = List{
	testOrd[3]: TestPatch{nil},
	testOrd[0]: TestPatch{nil},
	testOrd[1]: TestPatch{nil},
	testOrd[4]: TestPatch{nil},
	testOrd[2]: TestPatch{nil},
}

func TestErrors(t *testing.T) {
	a := assert.New(t)
	var err error

	patchList := List{
		"1.5.3": TestPatch{nil},
		"1.7.6": TestPatch{errNoReason},
	}

	err = patchList.Apply()
	if a.Error(err) {
		patchErr, ok := err.(*Error)
		if a.True(ok) {
			a.True(patchErr.IsApply)
			a.Equal("1.7.6", patchErr.InVersion.String())
			a.Equal(errNoReason, patchErr.OriginalError)
		}
	}

	err = patchList.Revert()
	if a.Error(err) {
		patchErr, ok := err.(*Error)
		if a.True(ok) {
			a.False(patchErr.IsApply)
			a.Equal("1.7.6", patchErr.InVersion.String())
			a.Equal(errNoReason, patchErr.OriginalError)
		}
	}
}

func TestHighestVersion(t *testing.T) {
	a := assert.New(t)

	ver, err := testPatchList.HighestVersion()
	if !a.Nil(err) {
		return
	}
	a.Equal("64.3.2", ver.String())
}

func TestPick(t *testing.T) {
	a := assert.New(t)

	c, _ := version.NewConstraint(">=1, <2")
	picked, err := testPatchList.Pick(c)
	if !a.Nil(err) {
		return
	}
	vers, err := picked.Versions()
	if a.Nil(err) {
		a.Equal(2, len(vers))
	}
}

func TestSortVersions(t *testing.T) {
	a := assert.New(t)

	vers, err := testPatchList.Versions()

	if !a.Nil(err) {
		return
	}

	if !a.Equal(len(testOrd), len(vers)) {
		return
	}

	for i, ver := range vers {
		a.Equal(testOrd[i], ver.String())
	}
}
