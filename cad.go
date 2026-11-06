// Code generated for the SimScale cadlib Go bindings. DO NOT EDIT.

package cadlib

/*
// cadlib.h lives in include/
//
// The native library is linked by name. Install libcadlib to a system
// location (e.g. /usr/local/lib, then ldconfig), or point cgo at a custom
// location via CGO_LDFLAGS="-L<dir>". See https://github.com/SimScaleGmbH/simscale-go-cadlib#install.
//
// The extra system libs (-lm -ldl -lpthread) are ELF/Linux-only; on macOS and
// Windows the C runtime provides those symbols, so link the cadlib only.
#cgo CFLAGS: -I${SRCDIR}/include
#cgo linux   LDFLAGS: -lcadlib -lm -ldl -lpthread
#cgo darwin  LDFLAGS: -lcadlib
#cgo windows LDFLAGS: -lcadlib
#include <stdlib.h>
#include <cadlib.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Error codes from the C library
const (
	CadLibOK               = 0
	CadLibErrorInvalidUUID = -1
	CadLibErrorNullPointer = -2
	CadLibErrorInvalidUTF8 = -3
	CadLibErrorInvalidJSON = -4
	CadLibErrorExpected    = -5
	CadLibErrorRuntime     = -6

	// Deprecated: use CadLibErrorExpected, which matches the C name CAD_ERROR_EXPECTED.
	CadLibErrorFailed = CadLibErrorExpected
)

// CadExpectedError is returned when a known, anticipated CAD error occurs.
// ErrorCode is the machine-readable identifier; ErrorDetails carries structured context.
type CadExpectedError struct {
	ErrorCode    string
	ErrorDetails map[string]any
}

func (e *CadExpectedError) Error() string {
	return fmt.Sprintf("Expected cad lib error: %s", e.ErrorCode)
}

// CadUnexpectedError is returned when an unexpected internal error occurs in the CAD library.
type CadUnexpectedError struct {
	Message string
}

func (e *CadUnexpectedError) Error() string {
	return e.Message
}

// CadErrorKind is the discriminant for the JSON error payload.
type CadErrorKind string

// Valid CadErrorKind values.
const (
	CadErrorKindExpected   CadErrorKind = "expected"
	CadErrorKindUnexpected CadErrorKind = "unexpected"
)

// cadErrorPayload mirrors the JSON structure.
type cadErrorPayload struct {
	Kind         CadErrorKind   `json:"kind"`
	ErrorCode    string         `json:"error_code"`
	ErrorDetails map[string]any `json:"error_details"`
	Message      string         `json:"message"`
}

// parseCadError decodes the JSON error payload from the FFI layer into a typed error.
// Falls back to CadUnexpectedError when the message is not valid JSON or has an unknown kind.
func parseCadError(rawMsg string) error {
	var payload cadErrorPayload
	if err := json.Unmarshal([]byte(rawMsg), &payload); err != nil {
		return &CadUnexpectedError{Message: rawMsg}
	}
	switch payload.Kind {
	case CadErrorKindExpected:
		return &CadExpectedError{
			ErrorCode:    payload.ErrorCode,
			ErrorDetails: payload.ErrorDetails,
		}
	case CadErrorKindUnexpected:
		return &CadUnexpectedError{Message: payload.Message}
	default:
		return &CadUnexpectedError{Message: rawMsg}
	}
}

// withErrorHandling calls fn with an error output pointer, frees the message on failure,
// and returns a typed error. Returns nil on success.
func withErrorHandling(fn func(**C.char) C.int) error {
	var errMsg *C.char
	if code := fn(&errMsg); code != CadLibOK {
		defer C.cad_free_string(errMsg)
		return parseCadError(C.GoString(errMsg))
	}
	return nil
}

// Cad represents a handle to a CAD model state.
type Cad struct {
	handle *C.CadHandle
}

// Load creates a new Cad instance from the CAD data map entry of a workflow method's run request,
// or from anything else marshalling to the same shape (a map, a struct of your own).
//
//	cad, err := cadlib.Load(request.DataMap.Cad)
//
// Only the entry's externalReference is required; dataType and source.type are validated when
// present and any other key is ignored. A caller holding raw JSON text passes it as a
// json.RawMessage, which marshals verbatim; a plain string would marshal to a JSON string and be
// rejected with invalid-cad-model-json. The error codes it can return are listed in
// https://github.com/SimScaleGmbH/simscale-go-cadlib/blob/main/ERRORS.md. An entry that
// cannot be marshalled to JSON at all returns a plain error.
func Load(entry any) (*Cad, error) {
	// Fail fast if the native library is incompatible with these bindings
	// (checked once per process; see version_check.go).
	if err := checkLibVersion(); err != nil {
		return nil, &CadUnexpectedError{Message: err.Error()}
	}

	cadModel, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal the cad data map entry: %w", err)
	}

	cCadModel := C.CString(string(cadModel))
	defer freeCString(cCadModel)

	var handle *C.CadHandle
	if err := withErrorHandling(func(e **C.char) C.int {
		return C.cad_load(cCadModel, &handle, e)
	}); err != nil {
		return nil, err
	}

	return &Cad{handle: handle}, nil
}

// Close releases the resources associated with this Cad instance.
// After calling Close, the Cad instance should not be used.
func (c *Cad) Close() {
	if c.handle != nil {
		C.cad_destroy(c.handle)
		c.handle = nil
	}
}

// Version returns the library version string.
func Version() string {
	cVersion := C.cad_version()
	defer C.cad_free_string(cVersion)
	return C.GoString(cVersion)
}

// executeQuery is a helper to execute a query and unmarshal the JSON result.
func (c *Cad) executeQuery(queryFn func(*C.char, **C.char, **C.char) C.int, params any, result any) error {
	if c.handle == nil {
		return errors.New("cad handle is nil or has been closed")
	}

	var paramsJSON *C.char
	// A nil *FooParameters reaches us as a non-nil interface holding a typed nil
	// pointer, so `params != nil` alone would marshal it to the JSON literal
	// "null". Queries without parameters must pass a null pointer.
	if params != nil && !reflect.ValueOf(params).IsNil() {
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal parameters: %w", err)
		}
		paramsJSON = C.CString(string(jsonBytes))
		defer freeCString(paramsJSON)
	}

	var resultJSON *C.char
	err := withErrorHandling(func(e **C.char) C.int {
		return queryFn(paramsJSON, &resultJSON, e)
	})
	if err != nil {
		return err
	}
	defer C.cad_free_string(resultJSON)

	if err := json.Unmarshal([]byte(C.GoString(resultJSON)), result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return nil
}

// ============================================================================
// Generated Query Methods
// ============================================================================

// Tessellate tessellates the model and writes the STL files into a directory.
//
// One STL file is written per region by default, or one per face when
// write_one_file_per_face is set. The directory is created if it does not
// exist, and files in it with the same names are overwritten.
//
// tessellation_type, max_facet_width and allowed_facet_faults_level
// apply to Parasolid models, and tessellation_type has to be set explicitly;
// feature_size, adaptivity and max_edge_length apply to nTop models
// and all default when omitted. tess_surface_area is set in the response
// only when calculate_area is requested.
//
// Supported for Parasolid, STL and nTop models.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - tessellation-has-error[-with-entity] - The kernel could not tessellate one or more
//     entities.
//   - tessellation-ntop-negative-adaptivity - adaptivity is negative.
//   - tessellation-ntop-zero-edge-length - max_edge_length is zero.
//   - tessellation-ntop-zero-feature-size - feature_size is zero.
//   - tessellation-too-large - The result could not be written; the kernel ran out of disk
//     space.
func (c *Cad) Tessellate(destDir string, params *TessellateParameters) (*TessellateResponse, error) {
	if c.handle == nil {
		return nil, errors.New("cad handle is nil or has been closed")
	}

	cDestDir := C.CString(destDir)
	defer freeCString(cDestDir)

	var paramsJSON *C.char
	if params != nil {
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameters: %w", err)
		}
		paramsJSON = C.CString(string(jsonBytes))
		defer freeCString(paramsJSON)
	}

	var resultJSON *C.char
	err := withErrorHandling(func(e **C.char) C.int {
		return C.cad_tessellate(c.handle, cDestDir, paramsJSON, &resultJSON, e)
	})
	if err != nil {
		return nil, err
	}
	defer C.cad_free_string(resultJSON)

	var result TessellateResponse
	if err := json.Unmarshal([]byte(C.GoString(resultJSON)), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

// TessellateFaces tessellates the requested faces, writing one STL file per face into a directory.
//
// Each requested face is tessellated into its own STL file, named after the
// face id; a face requested more than once is tessellated only once. The
// response additionally carries the normal of each face, evaluated at the
// mid point of the face.
//
// The faces to tessellate are required, so a query without parameters is
// rejected. The directory is created if it does not exist, and files in it with
// the same names are overwritten.
//
// Supported only for Parasolid models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - tessellation-has-error - The kernel could not tessellate one of the faces.
func (c *Cad) TessellateFaces(destDir string, params *TessellateFacesParameters) (*TessellateFacesResponse, error) {
	if c.handle == nil {
		return nil, errors.New("cad handle is nil or has been closed")
	}

	cDestDir := C.CString(destDir)
	defer freeCString(cDestDir)

	var paramsJSON *C.char
	if params != nil {
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameters: %w", err)
		}
		paramsJSON = C.CString(string(jsonBytes))
		defer freeCString(paramsJSON)
	}

	var resultJSON *C.char
	err := withErrorHandling(func(e **C.char) C.int {
		return C.cad_tessellate_faces(c.handle, cDestDir, paramsJSON, &resultJSON, e)
	})
	if err != nil {
		return nil, err
	}
	defer C.cad_free_string(resultJSON)

	var result TessellateFacesResponse
	if err := json.Unmarshal([]byte(C.GoString(resultJSON)), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

// GenerateComfortMapDisc generates the comfort-map disc and writes the STL result into a directory.
//
// The disc is a tessellated surface covering a region of interest: a circle of
// radius floor_radius around center, placed height above the ground.
// fineness controls how finely it is tessellated.
//
// By default the disc is written as a single file. Asking for one of the zipped
// encodings in result_encoding unpacks the result into the directory instead.
// The directory is created if it does not exist, and files in it with
// the same names are overwritten.
//
// A disc with no geometry is reported as empty-comfort-map-disc.
//
// Supported for Parasolid and STL models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - empty-comfort-map-disc - The query returned no geometry.
//   - lawson-cylinder-exceeds-size-box - floor_radius is too large for the model's size box.
//   - lawson-negative-height-above-ground - height is negative.
//   - lawson-roi-radius-must-be-positive - floor_radius is zero or negative.
//   - lawson-tessellation-max-width-must-be-positive - fineness is zero or negative.
//   - lawson-tessellation-too-fine - fineness produces more than 100 million triangles, which
//     the query does not allow.
//   - tessellation-has-error[-with-entity] - The kernel could not tessellate one or more
//     entities.
//   - tessellation-too-large - The result could not be written; the kernel ran out of disk
//     space.
func (c *Cad) GenerateComfortMapDisc(destDir string, params *GenerateComfortMapDiscParameters) (*GenerateComfortMapDiscResponse, error) {
	if c.handle == nil {
		return nil, errors.New("cad handle is nil or has been closed")
	}

	cDestDir := C.CString(destDir)
	defer freeCString(cDestDir)

	var paramsJSON *C.char
	if params != nil {
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameters: %w", err)
		}
		paramsJSON = C.CString(string(jsonBytes))
		defer freeCString(paramsJSON)
	}

	var resultJSON *C.char
	err := withErrorHandling(func(e **C.char) C.int {
		return C.cad_generate_comfort_map_disc(c.handle, cDestDir, paramsJSON, &resultJSON, e)
	})
	if err != nil {
		return nil, err
	}
	defer C.cad_free_string(resultJSON)

	var result GenerateComfortMapDiscResponse
	if err := json.Unmarshal([]byte(C.GoString(resultJSON)), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

// Volume computes the volume of each selected solid, and their total.
//
// The solids and sheets to measure are given by their entity ids, as
// entities returns them.
//
// volumes holds one entry per measured entity, keyed by its entity id,
// and total is their sum.
//
// Supported for Parasolid and nTop models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) Volume(params *VolumeParameters) (*VolumeResponse, error) {
	var result VolumeResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_volume(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BoundingBox computes the bounding box of the model, of each body, and of any faces given.
//
// The box of the whole model and one box per solid and sheet are always
// returned. Naming faces additionally returns one box per face, keyed by
// face id, and a single box enclosing all of them.
//
// Supported for Parasolid and nTop models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) BoundingBox(params *BoundingBoxParameters) (*BoundingBoxResponse, error) {
	var result BoundingBoxResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_bounding_box(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CenterOfGravity computes the center of gravity of the given faces and edges.
//
// The selection accepts faces and edges only; any other entity type fails
// the query. One center of gravity is computed for the selection as a
// whole, not one per selected entity.
//
// Supported only for Parasolid models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - compute-cog-entity-not-supported - A selected entity has unsupported type, it is neither
//     a face nor an edge.
func (c *Cad) CenterOfGravity(params *CenterOfGravityParameters) (*CenterOfGravityResponse, error) {
	var result CenterOfGravityResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_center_of_gravity(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DetectEnclosure finds the solids that enclose a given point.
//
// The response holds the id of every solid that contains it.
//
// Supported only for Parasolid models. Other model formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) DetectEnclosure(params *DetectEnclosureParameters) (*DetectEnclosureResponse, error) {
	var result DetectEnclosureResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_detect_enclosure(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DetectCommonRotation detects the rotation axis and center shared by the given faces.
//
// All the given faces must be rotationally aligned about a single axis; the
// query fails if they are not, or if no axis can be determined. The response
// gives the direction of that axis and a point on it.
//
// Supported only for Parasolid models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - detect-rotation-axis-not-consistent-for-all-faces - The selected faces do not share the
//     same rotation axis.
//   - detect-rotation-could-not-detect-axis - No rotation axis could be detected for the
//     selected faces.
//   - detect-rotation-flat-faces-selected[-with-entity] - The selection contains flat faces,
//     which have no rotation axis.
//   - detect-rotation-multiple-possible-directions - More than one rotation direction fits the
//     selected faces.
func (c *Cad) DetectCommonRotation(params *DetectCommonRotationParameters) (*DetectCommonRotationResponse, error) {
	var result DetectCommonRotationResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_detect_common_rotation(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Colors returns the colors defined on the model's solids and sheets.
//
// The query takes no parameters and covers the whole model. Each entry names
// the entity the color is defined on and the color itself, as a hex string
// such as #ff0000.
//
// Supported for Parasolid and nTop models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) Colors(params *ColorsParameters) (*ColorsResponse, error) {
	var result ColorsResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_colors(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MaterialPoints computes a point inside the material, for the whole model or for each solid.
//
// query_type chooses between the two. The response sets the field that
// matches it and leaves the other unset; the per-body field is keyed by
// entity id.
//
// A material point is only guaranteed to lie inside the material. Which
// interior point is returned is unspecified, and can differ between runs.
//
// Supported only for Parasolid models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) MaterialPoints(params *MaterialPointsParameters) (*MaterialPointsResponse, error) {
	var result MaterialPointsResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_material_points(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PlanarPatchPrimitives approximates groups of coplanar faces with a circle or rectangle primitive.
//
// Groups the requested faces into patches of adjacent, coplanar faces and
// approximates each patch with the primitive — a circle or a rectangle — whose
// area is closest to the area of the patch. All faces of a patch belong to the
// same occurrence. The query fails if any of the requested faces is not planar.
//
// Supported only for Parasolid models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
//
// # Query errors
//
//   - planar-patch-primitives-degenerate-patch[-with-entity] - A face has no interior point, or
//     a patch is not two-dimensional, so no primitive fits it.
//   - planar-patch-primitives-non-coplanar-faces[-with-entity] - Faces that are adjacent, and
//     so belong to the same patch, lie on different planes.
//   - planar-patch-primitives-non-planar-faces[-with-entity] - One or more of the requested
//     faces is not planar.
func (c *Cad) PlanarPatchPrimitives(params *PlanarPatchPrimitivesParameters) (*PlanarPatchPrimitivesResponse, error) {
	var result PlanarPatchPrimitivesResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_planar_patch_primitives(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ModelMetrics computes volume and surface-area metrics for the model and its entities.
//
// Computes aggregated volume and surface area for the model, along with
// per-solid, per-sheet, per-face and per-edge metrics.
//
// The query takes no parameters: the metrics always cover the whole model.
//
// Supported for Parasolid and nTop models. Other formats fail with
// query-type-unsupported-for-cad-format.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) ModelMetrics(params *ModelMetricsParameters) (*ModelMetricsResponse, error) {
	var result ModelMetricsResponse
	err := c.executeQuery(
		func(p *C.char, r **C.char, e **C.char) C.int {
			return C.cad_model_metrics(c.handle, p, r, e)
		},
		params,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Entities returns the model's entities: solids, sheets, faces, edges and vertices.
//
// The list is flat: each entity references its parents and children by their
// 0-based index in it. Solids and sheets are the top-level entities.
//
// Each entity carries the entity id that the other queries take as input.
//
// Supported for Parasolid, STL and nTop models.
//
// It returns a *CadExpectedError for a known domain error, or a
// *CadUnexpectedError if the library fails internally. Failing to
// marshal the parameters or unmarshal the result returns a plain error.
func (c *Cad) Entities() (*EntitiesResponse, error) {
	if c.handle == nil {
		return nil, errors.New("cad handle is nil or has been closed")
	}

	var resultJSON *C.char
	err := withErrorHandling(func(e **C.char) C.int {
		return C.cad_entities(c.handle, &resultJSON, e)
	})
	if err != nil {
		return nil, err
	}
	defer C.cad_free_string(resultJSON)

	var result EntitiesResponse
	if err := json.Unmarshal([]byte(C.GoString(resultJSON)), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}
