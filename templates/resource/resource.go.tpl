{{- template "boilerplate" }}

package {{ .Version }}

import (
    "github.com/anydotcloud/grm/pkg/compare"
    grmerr "github.com/anydotcloud/grm/pkg/error"
    "github.com/anydotcloud/grm/pkg/path/fieldpath"
    "github.com/anydotcloud/grm/pkg/types/resource"
    "github.com/anydotcloud/grm/pkg/types/resource/schema"

    resschema "{{ .ResourceSchemaPackage }}"
)

{{ .Documentation }}
type {{ .Kind.Name }} struct {
    values map[string]interface{}
    errors []error
}

// New returns a pointer to a new {{ .Kind.Name }}
func New() *{{ .Kind.Name }} {
    return &{{ .Kind.Name }}{
        values: map[string]interface{}{},
        errors: []error{},
    }
}

// IsValid returns true if the desired state has a possibility of being in
// sync with the latest observed state, false otherwise. If a user sets a
// desired field to a value that caused an Update operation to fail,
// IsValid() would return false for the resource and a call to Errors()
// would return a non-empty set of errors that placed the resource into an
// invalid state.
func (r *{{ .Kind.Name }}) IsValid() bool {
    return len(r.errors) == 0
}

// IsReady returns true if the resource's state indicates that the resource
// is "active", "available" or "ready".
func (r *{{ .Kind.Name }}) IsReady() bool {
    // TODO(jaypipes)
    return false
}

// IsImmutable returns true if the resource's state indicates that the resource
// may NOT be modified.
func (r *{{ .Kind.Name }}) IsImmutable() bool {
    // TODO(jaypipes)
    return false
}

// Errors returns zero or more errors that indicate why the resource may be
// in an invalid state.
func (r *{{ .Kind.Name }}) Errors() []error {
    return r.errors
}

// Identifiers returns an Identifiers which contain all the information
// needed to identify the resource.
func (r *{{ .Kind.Name }}) Identifiers() resource.Identifiers {
    return resschema.Identifiers,
}

// Schema returns a Schema that describes the resource's fields and
// identifiers
func (r *{{ .Kind.Name }}) Schema() schema.Schema {
    return resschema.Schema,
}

// Delta returns a Delta object containing the difference between this
// Resource and another.
func (r *{{ .Kind.Name }}) Delta(other resource.Resource) *compare.Delta {
    // TODO(jaypipes)
    return nil
}

// Values returns a map, keyed by stringified field path, of field
// values.
func (r *{{ .Kind.Name }}) Values() map[string]interface{} {
    return r.values
}

// ValueAt returns the value stored in the Resource for a Field identified by
// a supplied field path. Note that there is no way to retrieve a single
// element in a list field or a single key from a map field. Instead, Value
// returns the entire slice or map for a field identified by the supplied
// field path.
//
// Note that the field path is searched in a case-insensitive fashion. If there
// is no such field at the supplied path, returns a (nil, false) tuple.
func (r *{{ .Kind.Name }}) ValueAt(p *fieldpath.Path) (interface{}, bool) {
    for fp, v := range r.values {
        if strings.EqualFold(fp, p.String()) {
            return v, true
        }
    }
    return nil, false
}

// SetAt sets the value of a Resource field at the specified field path.
//
// Note that the field path is searched in a case-insensitive fashion. If there
// is no such field at the supplied path, returns an error.
func (r *{{ .Kind.Name }}) SetAt(p *fieldpath.Path, val interface{}) error {
    for fp := r.values {
        if strings.EqualFold(fp, p.String()) {
            r.values[fp] = val
            return nil
        }
    }
    return grmerr.UnknownFieldAtPath(p.String())
}
