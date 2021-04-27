# Patch

This library provides functionality for coordinating, applying, and reverting patches.

## Examples

This simple script demonstrates the syntax of a patch list and its simple application.

```go
package main

import (
	"fmt"

	"github.com/edge/patch"
	"github.com/hashicorp/go-version"
)

type DummyPatch struct {
	Message string
}

func (patch DummyPatch) Apply() error {
	fmt.Println(patch.Message)
	return nil
}

func (patch DummyPatch) Revert() error {
	fmt.Println(patch.Message)
	return nil
}

var patches = patch.List{
	"2.0.1": DummyPatch{"2.0 update inc prepatches"},
	"1.2.0": DummyPatch{"upgrade to 1.2 (rebalancing)"},
	"1.0.1": DummyPatch{"install day zero patch"},
	"1.0.0": DummyPatch{"1.0 install OK"},
}

func main() {
	c, _ := version.NewConstraint("*")
	if err := patches.Apply(c); err != nil {
		panic(err)
	}
}
```

Expected output:

```
1.0 install OK
install day zero patch
upgrade to 1.2 (rebalancing)
2.0 update inc prepatches
```

## Usage Notes

### Sorting Patches

It doesn't matter in which order you write your patch list - it will be automatically sorted based on the version numbers. Highest to lowest, or lowest to highest, either way is fine. (You could even organise them randomly, but that's probably not helpful.)

### Error Handling

Each patch may return an error, which you can handle in your own controller as appropriate based on your requirements. For example, you may want to add automatic rollbacks, or simply log errors, or provoke a user intervention. A patch list alone does not cover all possible scenarios - all it does is pick and execute all patches that match the given version constraint(s), and returns any error it encounters as soon as it happens.

When a patch produces an error while applying/reverting a patch list, a special [Error](./error.go) struct is returned providing additional information about the patch. Your controller can use this struct to decide how to handle the problem.

### Version Format

All versions specified in a patch list [must follow SemVer](http://semver.org/) - specifically, in the format that [go-version](https://github.com/hashicorp/go-version) expresses it. This ensures internal consistency and the ability to sort patches predictably.

If you specify an invalid patch version, you are likely to encounter errors in runtime, such as `unable to locate patch "xyz"` (reflecting that a version string, parsed through go-version back to a string, does not match itself).

If you want to organise patches in a different, non-numeric way (such as naming them) then this package may not be for you.

## Further Reading

`go doc` will tell you most of what you need to know!
