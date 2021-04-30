package patch

// Patch implementations apply and revert some effect.
//
// It is possible for a patch to return an error signifying a failure to apply or revert, or in some cases the impossibility of doing so.
// For example, a data migration from one schema to another may not be reversible.
type Patch interface {
	Apply() error  // Apply patch.
	Revert() error // Revert patch.
}
