// Package cadlib provides Go bindings for the SimScale cadlib native library.
//
// The package links against the prebuilt cadlib native library over its stable
// C ABI (libcadlib.so plus cadlib.h). Install the native library first, see
// https://github.com/SimScaleGmbH/simscale-go-cadlib#install, then load a model
// from the CAD data map entry of the method's run request and run queries against it:
//
//	cad, err := cadlib.Load(request.DataMap.Cad)
//	if err != nil {
//		return err
//	}
//	defer cad.Close()
//
//	resp, err := cad.Volume(&cadlib.VolumeParameters{EntityIDs: []string{"B1_TE5"}})
//
// Without the native library, building fails with "cannot find -lcadlib" and
// running with "error while loading shared libraries: libcadlib.so" (or the
// macOS and Windows equivalents). Install it with install-cadlib.sh, or
// install-cadlib.ps1 on Windows, from the release matching this module's
// version. For a non-default prefix, set CGO_LDFLAGS=-L<prefix>/lib to build
// and add that directory to the platform's library path to run.
//
// Every Cad value owns native resources, so Close must be called when the model
// is no longer needed. Query methods are blocking and are not safe for
// concurrent use on a single Cad value.
//
// Queries fail with one of two typed errors: [CadExpectedError] for known domain
// errors (carrying a machine-readable ErrorCode) and [CadUnexpectedError] for
// internal failures. Use errors.As to tell them apart. The one exception is a
// failure in the bindings themselves to marshal an argument or unmarshal a
// result to or from JSON, which is returned as a plain error.
//
// The model types in models.go and the query methods in cad.go are generated.
package cadlib
