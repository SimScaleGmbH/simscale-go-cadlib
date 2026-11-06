#ifndef CADLIB_H
#define CADLIB_H

/* AUTO-GENERATED - DO NOT EDIT */

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#ifdef __cplusplus
extern "C" {
#endif


/**
 * Result type for FFI queries.
 *
 * - `CAD_OK` (0): Success
 * - `CAD_ERROR_INVALID_UUID` (-1): Invalid UUID format
 * - `CAD_ERROR_NULL_POINTER` (-2): Null pointer passed
 * - `CAD_ERROR_INVALID_UTF8` (-3): Invalid UTF-8 string
 * - `CAD_ERROR_INVALID_JSON` (-4): Invalid JSON
 * - `CAD_ERROR_EXPECTED`     (-5): Expected error (CAD query failed)
 * - `CAD_ERROR_RUNTIME`      (-6): Runtime error
 */
#define CAD_OK 0

#define CAD_ERROR_INVALID_UUID -1

#define CAD_ERROR_NULL_POINTER -2

#define CAD_ERROR_INVALID_UTF8 -3

#define CAD_ERROR_INVALID_JSON -4

#define CAD_ERROR_EXPECTED -5

#define CAD_ERROR_RUNTIME -6

typedef enum CadEntityType {
    CAD_ENTITY_TYPE_SOLID,
    CAD_ENTITY_TYPE_SHEET,
    CAD_ENTITY_TYPE_FACE,
    CAD_ENTITY_TYPE_EDGE,
    CAD_ENTITY_TYPE_VERTEX,
} CadEntityType;

typedef enum CadResultEncoding {
    CAD_RESULT_ENCODING_ASCII,
    CAD_RESULT_ENCODING_BINARY,
    CAD_RESULT_ENCODING_ZIPPED_ASCII,
    CAD_RESULT_ENCODING_ZIPPED_BINARY,
} CadResultEncoding;

typedef enum CadMaterialPointsQueryType {
    CAD_MATERIAL_POINTS_QUERY_TYPE_ENTIRE_MODEL,
    CAD_MATERIAL_POINTS_QUERY_TYPE_PER_BODY,
} CadMaterialPointsQueryType;

typedef enum CadPatchType {
    CAD_PATCH_TYPE_RECTANGLE,
    CAD_PATCH_TYPE_CIRCLE,
} CadPatchType;

typedef enum CadTessellationType {
    CAD_TESSELLATION_TYPE_CHORD_TOL,
    CAD_TESSELLATION_TYPE_MAX_FACET_WIDTH,
} CadTessellationType;

typedef enum CadAllowedFacetFaultsLevel {
    CAD_ALLOWED_FACET_FAULTS_LEVEL_MILD,
    CAD_ALLOWED_FACET_FAULTS_LEVEL_MODERATE,
    CAD_ALLOWED_FACET_FAULTS_LEVEL_SERIOUS,
} CadAllowedFacetFaultsLevel;

/**
 * Handle to a Cad instance.
 */
typedef struct CadHandle {
    uint8_t private_[0];
} CadHandle;

typedef struct CadBox {
    double min_x;
    double min_y;
    double min_z;
    double max_x;
    double max_y;
    double max_z;
} CadBox;

typedef struct CadEntityBoxEntry {
    const char *key;
    struct CadBox value;
} CadEntityBoxEntry;

typedef struct CadBoundingBoxResponse {
    struct CadBox model_bounding_box;
    const struct CadEntityBoxEntry *bounding_box_per_body;
    uintptr_t bounding_box_per_body_len;
    const struct CadBox *faces_bounding_box;
    const struct CadEntityBoxEntry *bounding_box_per_face;
    uintptr_t bounding_box_per_face_len;
} CadBoundingBoxResponse;

typedef struct CadBoundingBoxParameters {
    const char *const *faces;
    uintptr_t faces_len;
} CadBoundingBoxParameters;

typedef struct CadVector3 {
    double x;
    double y;
    double z;
} CadVector3;

typedef struct CadCenterOfGravityResponse {
    struct CadVector3 cog;
} CadCenterOfGravityResponse;

typedef struct CadCenterOfGravityParameters {
    const char *const *entities;
    uintptr_t entities_len;
} CadCenterOfGravityParameters;

typedef struct CadColor {
    const char *entity_id;
    const char *color;
} CadColor;

typedef struct CadColorsResponse {
    const struct CadColor *items;
    uintptr_t items_len;
} CadColorsResponse;

typedef struct CadDetectCommonRotationResponse {
    struct CadVector3 center;
    struct CadVector3 axis;
} CadDetectCommonRotationResponse;

typedef struct CadDetectCommonRotationParameters {
    const char *const *faces;
    uintptr_t faces_len;
} CadDetectCommonRotationParameters;

typedef struct CadDetectEnclosureResponse {
    const char *const *items;
    uintptr_t items_len;
} CadDetectEnclosureResponse;

typedef struct CadDetectEnclosureParameters {
    struct CadVector3 point;
} CadDetectEnclosureParameters;

typedef struct CadEntity {
    const char *entity_id;
    const char *name;
    enum CadEntityType type;
    const int64_t *children;
    uintptr_t children_len;
    const int64_t *parents;
    uintptr_t parents_len;
} CadEntity;

typedef struct CadEntitiesResponse {
    const struct CadEntity *items;
    uintptr_t items_len;
} CadEntitiesResponse;

typedef struct CadGenerateComfortMapDiscResponse {
    int64_t number_of_triangles;
    const struct CadBox *bounding_box;
} CadGenerateComfortMapDiscResponse;

typedef struct CadVector2 {
    double x;
    double y;
} CadVector2;

typedef struct CadGenerateComfortMapDiscParameters {
    struct CadVector2 center;
    double height;
    double floor_radius;
    double fineness;
    const char *const *faces;
    uintptr_t faces_len;
    const double *ground_level;
    const enum CadResultEncoding *result_encoding;
} CadGenerateComfortMapDiscParameters;

typedef struct CadEntityVector3Entry {
    const char *key;
    struct CadVector3 value;
} CadEntityVector3Entry;

typedef struct CadMaterialPointsResponse {
    const struct CadVector3 *model_material_point;
    const struct CadEntityVector3Entry *material_point_per_body;
    uintptr_t material_point_per_body_len;
} CadMaterialPointsResponse;

typedef struct CadMaterialPointsParameters {
    enum CadMaterialPointsQueryType query_type;
} CadMaterialPointsParameters;

typedef struct CadVolumeMetric {
    const char *entity_id;
    double volume;
    double surface_area;
} CadVolumeMetric;

typedef struct CadFaceMetric {
    const char *entity_id;
    double area;
} CadFaceMetric;

typedef struct CadEdgeMetric {
    const char *entity_id;
    double length;
} CadEdgeMetric;

typedef struct CadModelMetricsResponse {
    double model_volume;
    double model_surface_area;
    const struct CadVolumeMetric *solids;
    uintptr_t solids_len;
    const struct CadVolumeMetric *sheets;
    uintptr_t sheets_len;
    const struct CadFaceMetric *faces;
    uintptr_t faces_len;
    const struct CadEdgeMetric *edges;
    uintptr_t edges_len;
} CadModelMetricsResponse;

typedef struct CadRectanglePrimitive {
    const struct CadVector3 *corners;
    uintptr_t corners_len;
    struct CadVector3 normal;
} CadRectanglePrimitive;

typedef struct CadCirclePrimitive {
    struct CadVector3 center;
    double radius;
    struct CadVector3 normal;
} CadCirclePrimitive;

typedef struct CadPlanarPatch {
    const char *const *faces;
    uintptr_t faces_len;
    enum CadPatchType patch_type;
    const struct CadRectanglePrimitive *rectangle;
    const struct CadCirclePrimitive *circle;
} CadPlanarPatch;

typedef struct CadPlanarPatchPrimitivesResponse {
    const struct CadPlanarPatch *patches;
    uintptr_t patches_len;
} CadPlanarPatchPrimitivesResponse;

typedef struct CadPlanarPatchPrimitivesParameters {
    const char *const *faces;
    uintptr_t faces_len;
} CadPlanarPatchPrimitivesParameters;

typedef struct CadTessellateFacesResponse {
    const struct CadEntityVector3Entry *face_normals;
    uintptr_t face_normals_len;
} CadTessellateFacesResponse;

typedef struct CadTessellateFacesParameters {
    const char *const *faces;
    uintptr_t faces_len;
    const enum CadTessellationType *tessellation_type;
    const double *max_facet_width;
    const enum CadAllowedFacetFaultsLevel *allowed_facet_faults_level;
} CadTessellateFacesParameters;

typedef struct CadTessellateResponse {
    int64_t num_bodies;
    const double *tess_surface_area;
} CadTessellateResponse;

typedef struct CadTessellateParameters {
    bool write_one_file_per_face;
    bool calculate_area;
    const enum CadTessellationType *tessellation_type;
    const double *max_facet_width;
    const enum CadAllowedFacetFaultsLevel *allowed_facet_faults_level;
    const double *feature_size;
    const double *adaptivity;
    const double *max_edge_length;
} CadTessellateParameters;

typedef struct CadEntityF64Entry {
    const char *key;
    double value;
} CadEntityF64Entry;

typedef struct CadVolumeResponse {
    double total;
    const struct CadEntityF64Entry *volumes;
    uintptr_t volumes_len;
} CadVolumeResponse;

typedef struct CadVolumeParameters {
    const char *const *entity_ids;
    uintptr_t entity_ids_len;
} CadVolumeParameters;

/**
 * Execute a `bounding_box` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_bounding_box(struct CadHandle *handle,
                     const char *parameters_json,
                     char **out_result,
                     char **out_error);

/**
 * Free a `CadBoundingBoxResponse` returned by `cad_bounding_box_t`.
 *
 * # Safety
 * `ptr` must come from `cad_bounding_box_t` or be null; do not use or free it twice.
 */
void cad_bounding_box_response_free(struct CadBoundingBoxResponse *ptr);

/**
 * Typed `bounding_box` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_bounding_box_response_free` and `*out_error` with `cad_free_string`
 */
int cad_bounding_box_t(struct CadHandle *handle,
                       const struct CadBoundingBoxParameters *params,
                       struct CadBoundingBoxResponse **out_result,
                       char **out_error);

/**
 * Execute a `center_of_gravity` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_center_of_gravity(struct CadHandle *handle,
                          const char *parameters_json,
                          char **out_result,
                          char **out_error);

/**
 * Free a `CadCenterOfGravityResponse` returned by `cad_center_of_gravity_t`.
 *
 * # Safety
 * `ptr` must come from `cad_center_of_gravity_t` or be null; do not use or free it twice.
 */
void cad_center_of_gravity_response_free(struct CadCenterOfGravityResponse *ptr);

/**
 * Typed `center_of_gravity` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_center_of_gravity_response_free` and `*out_error` with `cad_free_string`
 */
int cad_center_of_gravity_t(struct CadHandle *handle,
                            const struct CadCenterOfGravityParameters *params,
                            struct CadCenterOfGravityResponse **out_result,
                            char **out_error);

/**
 * Execute a `colors` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_colors(struct CadHandle *handle,
               const char *parameters_json,
               char **out_result,
               char **out_error);

/**
 * Free a `CadColorsResponse` returned by `cad_colors_t`.
 *
 * # Safety
 * `ptr` must come from `cad_colors_t` or be null; do not use or free it twice.
 */
void cad_colors_response_free(struct CadColorsResponse *ptr);

/**
 * Typed `colors` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_colors_response_free` and `*out_error` with `cad_free_string`
 */
int cad_colors_t(struct CadHandle *handle, struct CadColorsResponse **out_result, char **out_error);

/**
 * Free a Cad handle created by `cad_load`.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load` or null
 * - The handle must not be used after this call
 */
void cad_destroy(struct CadHandle *handle);

/**
 * Execute a `detect_common_rotation` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_detect_common_rotation(struct CadHandle *handle,
                               const char *parameters_json,
                               char **out_result,
                               char **out_error);

/**
 * Free a `CadDetectCommonRotationResponse` returned by `cad_detect_common_rotation_t`.
 *
 * # Safety
 * `ptr` must come from `cad_detect_common_rotation_t` or be null; do not use or free it twice.
 */
void cad_detect_common_rotation_response_free(struct CadDetectCommonRotationResponse *ptr);

/**
 * Typed `detect_common_rotation` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_detect_common_rotation_response_free` and `*out_error` with `cad_free_string`
 */
int cad_detect_common_rotation_t(struct CadHandle *handle,
                                 const struct CadDetectCommonRotationParameters *params,
                                 struct CadDetectCommonRotationResponse **out_result,
                                 char **out_error);

/**
 * Execute a `detect_enclosure` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_detect_enclosure(struct CadHandle *handle,
                         const char *parameters_json,
                         char **out_result,
                         char **out_error);

/**
 * Free a `CadDetectEnclosureResponse` returned by `cad_detect_enclosure_t`.
 *
 * # Safety
 * `ptr` must come from `cad_detect_enclosure_t` or be null; do not use or free it twice.
 */
void cad_detect_enclosure_response_free(struct CadDetectEnclosureResponse *ptr);

/**
 * Typed `detect_enclosure` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_detect_enclosure_response_free` and `*out_error` with `cad_free_string`
 */
int cad_detect_enclosure_t(struct CadHandle *handle,
                           const struct CadDetectEnclosureParameters *params,
                           struct CadDetectEnclosureResponse **out_result,
                           char **out_error);

/**
 * Get `entities` and return as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_entities(struct CadHandle *handle, char **out_result, char **out_error);

/**
 * Free a `CadEntitiesResponse` returned by `cad_entities_t`.
 *
 * # Safety
 * `ptr` must come from `cad_entities_t` or be null; do not use or free it twice.
 */
void cad_entities_response_free(struct CadEntitiesResponse *ptr);

/**
 * Typed `entities` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_entities_response_free` and `*out_error` with `cad_free_string`
 */
int cad_entities_t(struct CadHandle *handle,
                   struct CadEntitiesResponse **out_result,
                   char **out_error);

/**
 * Free a string allocated by the library.
 *
 * # Safety
 * - `s` must be a string allocated by this library or null
 */
void cad_free_string(char *s);

/**
 * Execute a `generate_comfort_map_disc` query, save results to `dest_dir`, and return metadata as JSON.
 *
 * # Parameters
 * - `handle`: Valid Cad handle
 * - `dest_dir`: Directory path to save results (must be valid UTF-8)
 * - `parameters_json`: Optional JSON string with `GenerateComfortMapDiscParameters` (null for defaults)
 * - `out_result`: Pointer to receive result JSON (caller must free with `cad_free_string`)
 * - `out_error`: Optional pointer to receive error message
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `dest_dir` must be a valid null-terminated UTF-8 string
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_generate_comfort_map_disc(struct CadHandle *handle,
                                  const char *dest_dir,
                                  const char *parameters_json,
                                  char **out_result,
                                  char **out_error);

/**
 * Free a `CadGenerateComfortMapDiscResponse` returned by `cad_generate_comfort_map_disc_t`.
 *
 * # Safety
 * `ptr` must come from `cad_generate_comfort_map_disc_t` or be null; do not use or free it twice.
 */
void cad_generate_comfort_map_disc_response_free(struct CadGenerateComfortMapDiscResponse *ptr);

/**
 * Typed `generate_comfort_map_disc` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_generate_comfort_map_disc_response_free` and `*out_error` with `cad_free_string`
 */
int cad_generate_comfort_map_disc_t(struct CadHandle *handle,
                                    const char *dest_dir,
                                    const struct CadGenerateComfortMapDiscParameters *params,
                                    struct CadGenerateComfortMapDiscResponse **out_result,
                                    char **out_error);

/**
 * Create a new Cad handle from the CAD data map entry of a workflow method's run request.
 *
 * # Parameters
 * - `cad_model`: A null-terminated C string holding the CAD data map entry as a JSON object, e.g.
 *   `{"externalReference":"cad-model.cad-id:[uuid]|cad-model.state-id:[uuid]"}`. Only
 *   `externalReference` is required; `dataType` and `source.type` are validated when present and
 *   any other key is ignored.
 * - `out_handle`: Pointer to receive the created handle
 * - `out_error`: Optional pointer to receive error message (caller must free with `cad_free_string`)
 *
 * # Returns
 * - `CAD_OK` on success
 * - Error code on failure
 *
 * # Errors
 * `CAD_ERROR_EXPECTED` carries one of these error codes (see
 * <https://github.com/SimScaleGmbH/simscale-go-cadlib/blob/main/ERRORS.md>):
 * - `invalid-cad-model-json`: The entry is not a JSON object.
 * - `unsupported-cad-model-data-type`: The entry declares a `dataType` this library is not compatible with.
 * - `unsupported-cad-model-source-type`: The entry declares a `source.type` other than `EXTERNAL`.
 * - `unsupported-cad-model-source`: The entry points at data on the local filesystem, which cannot be loaded.
 * - `unsupported-parametric-cad-model`: The entry holds parametric data, which cannot be loaded as a single CAD model.
 * - `invalid-cad-model`: The entry carries no CAD reference, or the reference is missing an id or holds a malformed UUID.
 *
 * # Safety
 * - `cad_model` must be a valid null-terminated string
 * - `out_handle` must be a valid pointer
 */
int cad_load(const char *cad_model,
             struct CadHandle **out_handle,
             char **out_error);

/**
 * Execute a `material_points` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_material_points(struct CadHandle *handle,
                        const char *parameters_json,
                        char **out_result,
                        char **out_error);

/**
 * Free a `CadMaterialPointsResponse` returned by `cad_material_points_t`.
 *
 * # Safety
 * `ptr` must come from `cad_material_points_t` or be null; do not use or free it twice.
 */
void cad_material_points_response_free(struct CadMaterialPointsResponse *ptr);

/**
 * Typed `material_points` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_material_points_response_free` and `*out_error` with `cad_free_string`
 */
int cad_material_points_t(struct CadHandle *handle,
                          const struct CadMaterialPointsParameters *params,
                          struct CadMaterialPointsResponse **out_result,
                          char **out_error);

/**
 * Execute a `model_metrics` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_model_metrics(struct CadHandle *handle,
                      const char *parameters_json,
                      char **out_result,
                      char **out_error);

/**
 * Free a `CadModelMetricsResponse` returned by `cad_model_metrics_t`.
 *
 * # Safety
 * `ptr` must come from `cad_model_metrics_t` or be null; do not use or free it twice.
 */
void cad_model_metrics_response_free(struct CadModelMetricsResponse *ptr);

/**
 * Typed `model_metrics` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_model_metrics_response_free` and `*out_error` with `cad_free_string`
 */
int cad_model_metrics_t(struct CadHandle *handle,
                        struct CadModelMetricsResponse **out_result,
                        char **out_error);

/**
 * Execute a `planar_patch_primitives` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_planar_patch_primitives(struct CadHandle *handle,
                                const char *parameters_json,
                                char **out_result,
                                char **out_error);

/**
 * Free a `CadPlanarPatchPrimitivesResponse` returned by `cad_planar_patch_primitives_t`.
 *
 * # Safety
 * `ptr` must come from `cad_planar_patch_primitives_t` or be null; do not use or free it twice.
 */
void cad_planar_patch_primitives_response_free(struct CadPlanarPatchPrimitivesResponse *ptr);

/**
 * Typed `planar_patch_primitives` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_planar_patch_primitives_response_free` and `*out_error` with `cad_free_string`
 */
int cad_planar_patch_primitives_t(struct CadHandle *handle,
                                  const struct CadPlanarPatchPrimitivesParameters *params,
                                  struct CadPlanarPatchPrimitivesResponse **out_result,
                                  char **out_error);

/**
 * Execute a `tessellate` query, save results to `dest_dir`, and return metadata as JSON.
 *
 * # Parameters
 * - `handle`: Valid Cad handle
 * - `dest_dir`: Directory path to save results (must be valid UTF-8)
 * - `parameters_json`: Optional JSON string with `TessellateParameters` (null for defaults)
 * - `out_result`: Pointer to receive result JSON (caller must free with `cad_free_string`)
 * - `out_error`: Optional pointer to receive error message
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `dest_dir` must be a valid null-terminated UTF-8 string
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_tessellate(struct CadHandle *handle,
                   const char *dest_dir,
                   const char *parameters_json,
                   char **out_result,
                   char **out_error);

/**
 * Execute a `tessellate_faces` query, save results to `dest_dir`, and return metadata as JSON.
 *
 * # Parameters
 * - `handle`: Valid Cad handle
 * - `dest_dir`: Directory path to save results (must be valid UTF-8)
 * - `parameters_json`: Optional JSON string with `TessellateFacesParameters` (null for defaults)
 * - `out_result`: Pointer to receive result JSON (caller must free with `cad_free_string`)
 * - `out_error`: Optional pointer to receive error message
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `dest_dir` must be a valid null-terminated UTF-8 string
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_tessellate_faces(struct CadHandle *handle,
                         const char *dest_dir,
                         const char *parameters_json,
                         char **out_result,
                         char **out_error);

/**
 * Free a `CadTessellateFacesResponse` returned by `cad_tessellate_faces_t`.
 *
 * # Safety
 * `ptr` must come from `cad_tessellate_faces_t` or be null; do not use or free it twice.
 */
void cad_tessellate_faces_response_free(struct CadTessellateFacesResponse *ptr);

/**
 * Typed `tessellate_faces` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_tessellate_faces_response_free` and `*out_error` with `cad_free_string`
 */
int cad_tessellate_faces_t(struct CadHandle *handle,
                           const char *dest_dir,
                           const struct CadTessellateFacesParameters *params,
                           struct CadTessellateFacesResponse **out_result,
                           char **out_error);

/**
 * Free a `CadTessellateResponse` returned by `cad_tessellate_t`.
 *
 * # Safety
 * `ptr` must come from `cad_tessellate_t` or be null; do not use or free it twice.
 */
void cad_tessellate_response_free(struct CadTessellateResponse *ptr);

/**
 * Typed `tessellate` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_tessellate_response_free` and `*out_error` with `cad_free_string`
 */
int cad_tessellate_t(struct CadHandle *handle,
                     const char *dest_dir,
                     const struct CadTessellateParameters *params,
                     struct CadTessellateResponse **out_result,
                     char **out_error);

/**
 * Get the library version string.
 *
 * # Returns
 * A pointer to a null-terminated string. The caller must free with `cad_free_string`.
 */
char *cad_version(void);

/**
 * Execute a `volume` query and return results as JSON.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `parameters_json` must be a valid null-terminated JSON string or null
 * - `out_result` and `out_error` must be valid pointers or null
 * - Caller must free `out_result` and `out_error` with `cad_free_string`
 */
int cad_volume(struct CadHandle *handle,
               const char *parameters_json,
               char **out_result,
               char **out_error);

/**
 * Free a `CadVolumeResponse` returned by `cad_volume_t`.
 *
 * # Safety
 * `ptr` must come from `cad_volume_t` or be null; do not use or free it twice.
 */
void cad_volume_response_free(struct CadVolumeResponse *ptr);

/**
 * Typed `volume` query.
 *
 * # Safety
 * - `handle` must be a valid handle from `cad_load`
 * - `out_result`/`out_error` must be valid pointers or null
 * - Free `*out_result` with `cad_volume_response_free` and `*out_error` with `cad_free_string`
 */
int cad_volume_t(struct CadHandle *handle,
                 const struct CadVolumeParameters *params,
                 struct CadVolumeResponse **out_result,
                 char **out_error);

#endif  /* CADLIB_H */

#ifdef __cplusplus
}
#endif
