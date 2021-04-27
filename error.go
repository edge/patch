package patch

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// Error represents a failure to apply or revert a patch.
// It is returned by a patch list only if a patch fails, so you can isolate patching errors:
//
//   if err := patches.Apply(); err != nil {
//     if patchErr, ok := err.(*patch.Error); ok {
//       // ...special handling for patch errors, e.g. rollback...
//     } else {
//       panic(err)
//     }
//   }
//
// Any other error returned in this package is just an ordinary Go error.
type Error struct {
	InVersion     *version.Version // Patch version in which the error occurred.
	IsApply       bool             // Whether the failed operation was apply or revert.
	OriginalError error            // The error encountered within the patch.
}

func (err *Error) Error() string {
	var verb string
	if err.IsApply {
		verb = "apply"
	} else {
		verb = "revert"
	}
	return fmt.Sprintf("failed to %s patch %s: %s", verb, err.InVersion, err.OriginalError)
}
