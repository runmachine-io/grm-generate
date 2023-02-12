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
	// Path is a "field path" that indicates where the field's value can be
	// found within the Resource.
	Path *fieldpath.Path
	// Config contains the configuration options for this field
	Config *config.FieldConfig
	// Definition contains metadata about the field's type
	Definition *FieldDefinition
}

// Names returns the set of normalized name variations for the field
func (f *Field) Names() names.Names {
	return names.New(f.Path.Back())
}

// NewField returns an initialized Field from a field path, configuration and
// FieldDefinition. We normalize each part of the supplied field path, so for
// example, "RegistryId" becomes "RegistryID" and "EncryptionConfig.KmsKeyId"
// becomes "EncryptionConfig.KMSKeyID".
func NewField(
	path *fieldpath.Path,
	cfg *config.FieldConfig,
	def *FieldDefinition,
) *Field {
	normPath := &fieldpath.Path{}
	for {
		part := path.PopFront()
		if part == "" {
			break
		}
		normed := names.New(part)
		normPath.PushBack(normed.Camel)
	}
	return &Field{
		Path:       normPath,
		Config:     cfg,
		Definition: def,
	}
}
