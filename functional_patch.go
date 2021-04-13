package patch

// FunctionalPatch is a simple Patch that forwards Apply/Revert calls to user functions.
// For reasons of structural clarity you should use concrete types in user code.
// This type is provided only for testing and advanced usage.
type FunctionalPatch struct {
	ApplyFunc  func() error
	RevertFunc func() error
}

// Apply patch.
func (p FunctionalPatch) Apply() error {
	return p.ApplyFunc()
}

// Revert patch.
func (p FunctionalPatch) Revert() error {
	return p.RevertFunc()
}
