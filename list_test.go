package patch

import (
	"testing"

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

var testOrd []string = []string{"0.1.0", "1.0.0", "1.7.3", "2.0.1", "64.3.2"}

// versions are jumbled in this test to ensure failure case if patch list doesn't sort correctly.
// testPatchListease don't do this in actual code!
var testPatchList List = List{
	testOrd[3]: TestPatch{nil},
	testOrd[0]: TestPatch{nil},
	testOrd[1]: TestPatch{nil},
	testOrd[4]: TestPatch{nil},
	testOrd[2]: TestPatch{nil},
}

func TestHighestVersion(t *testing.T) {
	a := assert.New(t)

	ver, err := testPatchList.HighestVersion()
	if !a.Nil(err) {
		return
	}
	a.Equal("64.3.2", ver.String())
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
