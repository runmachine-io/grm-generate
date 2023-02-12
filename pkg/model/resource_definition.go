// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package model

import (
	"sort"
	"strings"

	"github.com/samber/lo"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm/pkg/path/fieldpath"
)

// ResourceDefinition describes a single top-level resource in a cloud service
// API
type ResourceDefinition struct {
	// Config contains the resource-specific configuration options
	Config *config.ResourceConfig
	// Kind is the type of Resource
	Kind Kind
	// fields is a map, keyed by the **field path**, of Field objects
	// representing a field in the Resource.
	fields map[string]*Field
}

// FieldPaths returns a sorted list of field paths for this resource.
func (d *ResourceDefinition) GetFieldPaths() []*fieldpath.Path {
	pathStrs := lo.Keys(d.fields)
	sort.Strings(pathStrs)
	res := make([]*fieldpath.Path, len(pathStrs))
	for x, pathStr := range pathStrs {
		res[x] = fieldpath.FromString(pathStr)
	}
	return res
}

// GetField returns a Field given a field path. The search is case-insensitive
func (d *ResourceDefinition) GetField(path *fieldpath.Path) *Field {
	for pathStr, f := range d.fields {
		if strings.EqualFold(path.String(), pathStr) {
			return f
		}
	}
	return nil
}

// AddField adds a new Field to the resource definition at the supplied field
// path
func (d *ResourceDefinition) AddField(f *Field) {
	d.fields[f.Path.String()] = f
}

// NewResourceDefinition returns a pointer to a new ResourceDefinition that
// describes a single top-level resource in a cloud service API. Add fields to
// the ResourceDefinition by calling the AddField method.
func NewResourceDefinition(
	cfg *config.ResourceConfig,
	kind Kind,
) *ResourceDefinition {
	return &ResourceDefinition{
		Config: cfg,
		Kind:   kind,
		fields: map[string]*Field{},
	}
}
