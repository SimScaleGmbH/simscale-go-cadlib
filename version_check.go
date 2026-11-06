// The Go bindings accept exactly one native libcadlib: the one released with
// them.
//
// checkLibVersion compares the generated libVersion against the native lib's
// cad_version() and fails on any difference; the error points at the install
// scripts of the expected release.

package cadlib

import (
	"fmt"
	"sync"
)

var (
	versionCheckOnce sync.Once
	versionCheckErr  error
)

func checkLibVersion() error {
	versionCheckOnce.Do(func() {
		versionCheckErr = ensureLibVersion(libVersion, Version())
	})
	return versionCheckErr
}

// ensureLibVersion reports an error unless the native lib version is exactly expected.
func ensureLibVersion(expected, lib string) error {
	if lib == expected {
		return nil
	}
	return fmt.Errorf(
		"native libcadlib v%s does not match these Go bindings (v%s); install it with install-cadlib.sh or install-cadlib.ps1 from the v%s release",
		lib, expected, expected,
	)
}
