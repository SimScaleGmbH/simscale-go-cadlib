package cadlib

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"

// freeCString releases a string allocated with C.CString.
func freeCString(p *C.char) {
	C.free(unsafe.Pointer(p)) // nosemgrep: go.lang.security.audit.unsafe.use-of-unsafe-block
}
