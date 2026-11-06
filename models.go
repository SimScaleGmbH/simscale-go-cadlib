// Code generated for the SimScale cadlib Go bindings. DO NOT EDIT.
package cadlib

// AllowedFacetFaultsLevel level of faults allowed in tessellation (Supported only for Parasolid models)
// Used in quality checks to classify the tessellation faults in three levels: mild, moderate and serious.
// The default is `Mild`, which allows only minor faults that do not significantly affect the quality of the tessellation.
// `Moderate` allows more significant faults that may affect the quality but are still acceptable for most applications.
// `Serious` allows all faults, including those that may result in a poor quality tessellation, and should be used with caution.
type AllowedFacetFaultsLevel string

// Valid AllowedFacetFaultsLevel values.
const (
	AllowedFacetFaultsLevelMild     AllowedFacetFaultsLevel = "mild"
	AllowedFacetFaultsLevelModerate AllowedFacetFaultsLevel = "moderate"
	AllowedFacetFaultsLevelSerious  AllowedFacetFaultsLevel = "serious"
)

// BoundingBoxParameters parameters for bounding box info query
type BoundingBoxParameters struct {
	// Ids of faces to measure in addition to the model and its bodies. When
	// given, the response also holds the box enclosing all of them and one box
	// per face
	Faces []string `json:"faces,omitempty"`
}

// BoundingBoxResponse response for bounding box query
type BoundingBoxResponse struct {
	// Bounding box of the entire model
	ModelBoundingBox Box `json:"modelBoundingBox"`
	// Bounding box of each solid and sheet in the model, keyed by its entity id
	BoundingBoxPerBody map[string]Box `json:"boundingBoxPerBody"`
	// Single box enclosing all the requested faces. Set only when faces were
	// requested
	FacesBoundingBox *Box `json:"facesBoundingBox,omitempty"`
	// Bounding box of each requested face, keyed by its face id. Set only when
	// faces were requested
	BoundingBoxPerFace map[string]Box `json:"boundingBoxPerFace,omitempty"`
}

// Box 3D bounding box representation
type Box struct {
	// Minimum X coordinate
	MinX float64 `json:"minX"`
	// Minimum Y coordinate
	MinY float64 `json:"minY"`
	// Minimum Z coordinate
	MinZ float64 `json:"minZ"`
	// Maximum X coordinate
	MaxX float64 `json:"maxX"`
	// Maximum Y coordinate
	MaxY float64 `json:"maxY"`
	// Maximum Z coordinate
	MaxZ float64 `json:"maxZ"`
}

// CenterOfGravityParameters parameters for compute center of gravity query
type CenterOfGravityParameters struct {
	// Ids of the faces and edges to combine into a single center of gravity, as
	// `entities` returns them. Any other entity type fails the query.
	// Must not be empty
	Entities []string `json:"entities"`
}

// CenterOfGravityResponse response for compute center of gravity query
type CenterOfGravityResponse struct {
	// Center of gravity as a 3D vector
	Cog Vector3 `json:"cog"`
}

// CirclePrimitive circle approximating a planar patch
type CirclePrimitive struct {
	// Center of the circle
	Center Vector3 `json:"center"`
	// Radius of the circle
	Radius float64 `json:"radius"`
	// Outward normal of the patch
	Normal Vector3 `json:"normal"`
}

// Color information for an entity
type Color struct {
	// Id of the solid or sheet the color is defined on
	EntityID string `json:"entity_id"`
	// Color in hex format (e.g., "#ff0000")
	Color string `json:"color"`
}

// ColorsParameters parameters for get model colors query (empty)
type ColorsParameters struct{}

// MarshalJSON encodes the parameter-less query as JSON null, which is what
// the library expects for this schema (`"type": "null"`).
func (ColorsParameters) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// UnmarshalJSON accepts the JSON null this type marshals to.
func (*ColorsParameters) UnmarshalJSON([]byte) error {
	return nil
}

// ColorsResponse response for get model colors query
// List of colors defined on solids and sheets
type ColorsResponse []Color

// DetectCommonRotationParameters parameters for detect rotation query
type DetectCommonRotationParameters struct {
	// Ids of the faces to check for common rotation, as `entities` returns them. The query will check if these faces are rotationally aligned and share the same rotation center and axis. Must not be empty
	Faces []string `json:"faces"`
}

// DetectCommonRotationResponse response for detect common rotation query
type DetectCommonRotationResponse struct {
	// Center of common rotation
	Center Vector3 `json:"center"`
	// Axis of common rotation
	Axis Vector3 `json:"axis"`
}

// DetectEnclosureParameters parameters for detect enclosure query
type DetectEnclosureParameters struct {
	// Point to check for enclosure
	Point Vector3 `json:"point"`
}

// DetectEnclosureResponse response for detect enclosure query
// List of solids enclosing the input point
type DetectEnclosureResponse []string

// EdgeMetric length metric for a single edge
type EdgeMetric struct {
	// Id of the edge
	EntityID string `json:"entityId"`
	// Length of the edge
	Length float64 `json:"length"`
}

// EntitiesResponse response for entities query — a flat list of all entities.
//
// Each entity references its parents and children by their 0-based index in the list.
// Solids and sheets are the top-level entities.
type EntitiesResponse []Entity

// Entity an entity in the CAD model topology (solids, sheets, faces, edges, vertices).
type Entity struct {
	// Unique entity id, e.g. `A1_I2_B3_TE4` or `B1_TE2`
	EntityID string `json:"entityId"`
	// Human-readable entity name, e.g. `face_name@body_name`
	Name string `json:"name"`
	// Type of the entity
	Type EntityType `json:"type"`
	// Indices of the children of this entity in this list
	Children []int32 `json:"children"`
	// Indices in this list of all entities this entity is a child of.
	// Empty for top-level entities, i.e. solids and sheets.
	// An entity can have several parents: an edge is usually shared by two faces,
	// a vertex by several edges.
	Parents []int32 `json:"parents"`
}

// EntityType type of entity. A body is never part of the topology, since it is
// represented by its top-level solids or sheets.
type EntityType string

// Valid EntityType values.
const (
	EntityTypeSolid  EntityType = "SOLID"
	EntityTypeSheet  EntityType = "SHEET"
	EntityTypeFace   EntityType = "FACE"
	EntityTypeEdge   EntityType = "EDGE"
	EntityTypeVertex EntityType = "VERTEX"
)

// FaceMetric surface area metric for a single face
type FaceMetric struct {
	// Id of the face
	EntityID string `json:"entityId"`
	// Surface area of the face
	Area float64 `json:"area"`
}

// GenerateComfortMapDiscParameters parameters for comfort maps query
type GenerateComfortMapDiscParameters struct {
	// Center point of the region of interest on the ground plane. The disc will be generated in a circular region around this center point.
	Center Vector2 `json:"center"`
	// Height above ground of the region of interest (distance between the center point and the plane where the disc will be generated). Must not be negative
	Height float64 `json:"height"`
	// Radius of the region of interest on the ground plane (disc will be generated in a circular region with this radius around the center point). Must be greater than zero
	FloorRadius float64 `json:"floor_radius"`
	// Tessellation fineness parameter to control the size of the disc's tessellation,. Must be greater than zero, and large enough that the disc stays under the limit of 100 million triangles
	Fineness float64 `json:"fineness"`
	// Ids of the faces that make up the ground (optional)
	Faces []string `json:"faces,omitempty"`
	// Height of the ground plane. Defaults to 0 when omitted
	GroundLevel *float64 `json:"ground_level,omitempty"`
	// Defines the encoding of the file which contains the disc (optional). With `ascii` or `binary`, or when omitted, the disc is written as a single file; with `zipped-ascii` or `zipped-binary` the result is unpacked into the directory
	ResultEncoding *ResultEncoding `json:"result_encoding,omitempty"`
}

// GenerateComfortMapDiscResponse response for comfort maps query
type GenerateComfortMapDiscResponse struct {
	// Number of triangles in the tessellated disc. Always greater than zero: a
	// disc with no geometry fails the query instead
	NumberOfTriangles int32 `json:"number_of_triangles"`
	// Bounding box of the tessellated disc (optional)
	BoundingBox *Box `json:"bounding_box,omitempty"`
}

// MaterialPointsParameters parameters for material points query
type MaterialPointsParameters struct {
	// Defines whether one point is computed for the whole model or one per solid
	// and sheet. The response sets the matching field and leaves the other unset
	QueryType MaterialPointsQueryType `json:"queryType"`
}

// MaterialPointsQueryType query type for material points
type MaterialPointsQueryType string

// Valid MaterialPointsQueryType values.
const (
	MaterialPointsQueryTypeEntireModel MaterialPointsQueryType = "entireModel"
	MaterialPointsQueryTypePerBody     MaterialPointsQueryType = "perBody"
)

// MaterialPointsResponse response for material points query
type MaterialPointsResponse struct {
	// A point inside the model's material. Set when `queryType` is `entireModel`
	ModelMaterialPoint *Vector3 `json:"modelMaterialPoint,omitempty"`
	// A point inside each solid and sheet, keyed by its entity id. Set when
	// `queryType` is `perBody`
	MaterialPointPerBody map[string]Vector3 `json:"materialPointPerBody,omitempty"`
}

// ModelMetricsParameters parameters for model metrics query
type ModelMetricsParameters struct{}

// ModelMetricsResponse response for model metrics query
type ModelMetricsResponse struct {
	// Sum of the volumes in `solids`
	ModelVolume float64 `json:"modelVolume"`
	// Aggregated surface area of all solids and sheets
	ModelSurfaceArea float64 `json:"modelSurfaceArea"`
	// Volume and surface area metrics per solid
	Solids []VolumeMetric `json:"solids"`
	// Volume and surface area metrics per sheet
	Sheets []VolumeMetric `json:"sheets"`
	// Surface area metrics per face
	Faces []FaceMetric `json:"faces"`
	// Length metrics per edge
	Edges []EdgeMetric `json:"edges"`
}

// PatchType the kind of primitive that approximates a planar patch
type PatchType string

// Valid PatchType values.
const (
	PatchTypeRectangle PatchType = "rectangle"
	PatchTypeCircle    PatchType = "circle"
)

// PlanarPatch one patch of adjacent, coplanar faces, together with the primitive approximating it.
//
// `patch_type` names which primitive describes the patch: exactly one of `rectangle` and
// `circle` is set, the one `patch_type` names.
type PlanarPatch struct {
	// The faces that make up this patch
	Faces []string `json:"faces"`
	// The kind of primitive that approximates this patch
	PatchType PatchType `json:"patchType"`
	// The rectangle approximating this patch, set when `patch_type` is `rectangle`
	Rectangle *RectanglePrimitive `json:"rectangle,omitempty"`
	// The circle approximating this patch, set when `patch_type` is `circle`
	Circle *CirclePrimitive `json:"circle,omitempty"`
}

// PlanarPatchPrimitivesParameters parameters for planar patch primitives query
type PlanarPatchPrimitivesParameters struct {
	// Ids of the planar faces to group into patches, as `entities` returns them. The query fails if any of the given faces is not planar. Must not be empty
	Faces []string `json:"faces"`
}

// PlanarPatchPrimitivesResponse response for planar patch primitives query
type PlanarPatchPrimitivesResponse struct {
	// The patches found among the requested faces
	Patches []PlanarPatch `json:"patches"`
}

// RectanglePrimitive rectangle approximating a planar patch
type RectanglePrimitive struct {
	// The four corners of the rectangle, counterclockwise around `normal`
	Corners []Vector3 `json:"corners"`
	// Outward normal of the patch
	Normal Vector3 `json:"normal"`
}

// ResultEncoding encoding format for STL result data, used in some queries to control how the result is encoded (e.g. for comfort map disc)
type ResultEncoding string

// Valid ResultEncoding values.
const (
	ResultEncodingASCII        ResultEncoding = "ascii"
	ResultEncodingBinary       ResultEncoding = "binary"
	ResultEncodingZippedASCII  ResultEncoding = "zipped-ascii"
	ResultEncodingZippedBinary ResultEncoding = "zipped-binary"
)

// TessellateFacesParameters parameters for the tessellate faces query.
//
// Tessellates a chosen set of faces: each face becomes its own STL file, named
// after the face id. A face listed more than once is tessellated only once.
//
// Supported only for Parasolid models.
type TessellateFacesParameters struct {
	// Ids of the faces to tessellate, as `entities` returns them. Each face is
	// tessellated into its own STL file, named after the face id. A face
	// requested more than once is tessellated only once. Must not be empty
	Faces []string `json:"faces"`
	// Defines the tessellation strategy. Defaults to `chordTol` when omitted
	TessellationType *TessellationType `json:"tessellationType,omitempty"`
	// Maximum facet width for tessellationType=maxFacetWidth. If omitted, it is
	// derived from the size of the face
	MaxFacetWidth *float64 `json:"maxFacetWidth,omitempty"`
	// Level of faults allowed in the tessellation. Used in quality checks.
	// Defaults to `mild` when omitted
	AllowedFacetFaultsLevel *AllowedFacetFaultsLevel `json:"allowedFacetFaultsLevel,omitempty"`
}

// TessellateFacesResponse response for the tessellate faces query
type TessellateFacesResponse struct {
	// Unit normal of each requested face, evaluated at the mid point of the
	// face, keyed by face id
	FaceNormals map[string]Vector3 `json:"faceNormals"`
}

// TessellateParameters parameters for tessellate query
type TessellateParameters struct {
	// If true, one STL file is written per face; if false, one per region
	WriteOneFilePerFace bool `json:"writeOneFilePerFace"`
	// If true, computes the surface area of the tessellation and returns it as
	// `tess_surface_area` in the response
	CalculateArea bool `json:"calculateArea"`
	// Supported only for Parasolid models. Defines the tessellation strategy.
	// Has no default: the query fails with `invalid-parameter` if it is not set
	TessellationType *TessellationType `json:"tessellationType,omitempty"`
	// Supported only for Parasolid models. Maximum facet width, used when `tessellation_type`
	// is `maxFacetWidth`. If omitted, it is computed from the bounding box of the body
	MaxFacetWidth *float64 `json:"maxFacetWidth,omitempty"`
	// Supported only for Parasolid models. Level of faults allowed in the tessellation.
	// Used in quality checks. Defaults to `mild` when omitted
	AllowedFacetFaultsLevel *AllowedFacetFaultsLevel `json:"allowedFacetFaultsLevel,omitempty"`
	// Supported only for nTop models. Granularity at which the solid's field is sampled.
	// Make this smaller to get more triangles and lower max edge length.
	// Must be greater than zero. Defaults to 0.0015 when omitted
	FeatureSize *float64 `json:"featureSize,omitempty"`
	// Supported only for nTop models. Controls how much the mesh is decimated.
	// Make this higher to get less triangles and higher max edge length.
	// Must be greater than zero. Defaults to 10 when omitted
	Adaptivity *float64 `json:"adaptivity,omitempty"`
	// Supported only for nTop models. Specifies the maximum allowable length for any edge in the tessellated
	// mesh. During refinement, edges longer than this value are subdivided to ensure the
	// mesh respects this limit. Must be greater than zero. Defaults to 10 when omitted
	MaxEdgeLength *float64 `json:"maxEdgeLength,omitempty"`
}

// TessellateResponse response for tessellation query
type TessellateResponse struct {
	// Number of bodies in the tessellated geometry
	NumBodies int32 `json:"numBodies"`
	// Surface area of the tessellated geometry. Set only when `calculate_area`
	// was requested
	TessSurfaceArea *float64 `json:"tessSurfaceArea,omitempty"`
}

// TessellationType tessellation strategy (Supported only for Parasolid models)
// `ChordTol` controls curve tolerance by specifying the maximum allowable distance between a true curved edge
// and the straight-line chords (facet edges) used to approximate it.
// The `MaxFacetWidth` strategy uses a maximum facet width to control the maximum size of the facets, which can be more intuitive in some cases.
// When using `MaxFacetWidth`, the `maxFacetWidth` parameter controls the maximum width of the facets.
type TessellationType string

// Valid TessellationType values.
const (
	TessellationTypeChordTol      TessellationType = "chordTol"
	TessellationTypeMaxFacetWidth TessellationType = "maxFacetWidth"
)

// Vector2 2D vector representation
type Vector2 struct {
	// X coordinate
	X float64 `json:"x"`
	// Y coordinate
	Y float64 `json:"y"`
}

// Vector3 3D vector representation
type Vector3 struct {
	// X coordinate
	X float64 `json:"x"`
	// Y coordinate
	Y float64 `json:"y"`
	// Z coordinate
	Z float64 `json:"z"`
}

// VolumeMetric volume and surface area metrics for a single solid or sheet
type VolumeMetric struct {
	// Id of the solid or sheet
	EntityID string `json:"entityId"`
	// Volume of the entity, when applicable
	Volume float64 `json:"volume"`
	// Surface area of the entity
	SurfaceArea float64 `json:"surfaceArea"`
}

// VolumeParameters parameters for compute volume query
type VolumeParameters struct {
	// Ids of the solids and sheets to measure, as `entities` returns them.
	// Must not be empty
	EntityIDs []string `json:"entityIds"`
}

// VolumeResponse response for volume query
type VolumeResponse struct {
	// Sum of the volumes in `volumes`
	Total float64 `json:"total"`
	// Volume of each measured entity, keyed by its entity id
	Volumes map[string]float64 `json:"volumes"`
}
