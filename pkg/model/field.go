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
	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm/pkg/names"
	"github.com/anydotcloud/grm/pkg/path/fieldpath"
)

// Field represents a single field in the Resource's Schema.
type Field struct {
	// Names is a set of normalized name variations for the field
	Names names.Names
	// Path is a "field path" that indicates where the field's value can be
	// found within the Resource.
	Path *fieldpath.Path
	// Config contains the configuration options for this field
	Config *config.FieldConfig
	// Definition contains metadata about the field's type
	Definition *FieldDefinition
}

// NewField returns an initialized Field
func NewField(
	names names.Names,
	path *fieldpath.Path,
	cfg *config.FieldConfig,
	def *FieldDefinition,
) *Field {
	return &Field{
		Names:      names,
		Path:       path,
		Config:     cfg,
		Definition: def,
	}
}
