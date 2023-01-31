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

	"github.com/samber/lo"

	"github.com/anydotcloud/grm-generate/pkg/config"
)

// ResourceDefinition describes a single top-level resource in a cloud service
// API
type ResourceDefinition struct {
	// Config contains the resource-specific configuration options
	Config *config.ResourceConfig
	// Kind is the type of Resource
	Kind Kind
	// Fields is a map, keyed by the **field path**, of Field objects
	// representing a field in the Resource.
	Fields map[string]*Field
}

// FieldPaths returns a sorted list of field paths for this resource.
func (d *ResourceDefinition) FieldPaths() []string {
	paths := lo.Keys(d.Fields)
	sort.Strings(paths)
	return paths
}

// NewResourceDefinition returns a pointer to a new ResourceDefinition that
// describes a single top-level resource in a cloud service API
func NewResourceDefinition(
	cfg *config.ResourceConfig,
	kind Kind,
	fields map[string]*Field, // map of fields, keyed by **field path**
) *ResourceDefinition {
	return &ResourceDefinition{
		Config: cfg,
		Kind:   kind,
		Fields: fields,
	}
}
