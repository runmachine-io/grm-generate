{{- template "boilerplate" }}

package schema

import (
    "strings"

    "github.com/anydotcloud/grm/pkg/path/fieldpath"
	"github.com/anydotcloud/grm/pkg/types/resource"
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
)

var (
    schemaFields = map[string]schema.Field{
{{- range $schemaFieldPathString, $schemaFieldTypeName := .Fields }}
        "{{ $schemaFieldPathString }}": {{ $schemaFieldTypeName }},
{{ end -}}
    }
)

type schema struct{
    schema.Kind
}

// Field returns a Field at a given field path, or nil if there is no Field
// at that path.
func (s *schema) Field(p *fieldpath.Path) schema.Field {
    for pathStr, f := range schemaFields {
        if strings.EqualFold(pathStr, p.String()) {
            return f
        }
    }
    return nil
}

// Fields returns a map, keyed by field path string, of Fields that
// describe the resource's member fields.
func (s *schema) Fields() map[string]schema.Field {
    return schemaFields
}

// Identifiers returns information about a resource's identifying fields
// and those fields' values.
Identifiers() schema.Identifiers {
    return Identifiers
}

// Schema contains methods that returns information about a resource's schema.
var Schema schema.Schema = &schema{Kind}
