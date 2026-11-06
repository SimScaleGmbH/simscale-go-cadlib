// These tests call the real native lib, so they need it on the cgo link and
// runtime library paths.

package cadlib

import "testing"

func TestLinkedLibMatchesBindings(t *testing.T) {
	if err := checkLibVersion(); err != nil {
		t.Fatalf("checkLibVersion() = %v, want nil", err)
	}
}

func TestMismatchedLibIsRejected(t *testing.T) {
	lib := Version()
	expected := libVersion + "-mismatch"
	want := "native libcadlib v" + lib + " does not match these Go bindings (v" + expected +
		"); install it with install-cadlib.sh or install-cadlib.ps1 from the v" + expected +
		" release"
	err := ensureLibVersion(expected, lib)
	if err == nil || err.Error() != want {
		t.Fatalf("ensureLibVersion(%q, %q) = %v, want %q", expected, lib, err, want)
	}
}
