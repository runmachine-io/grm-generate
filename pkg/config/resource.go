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

package config

import (
	"strings"

	"github.com/anydotcloud/grm/pkg/path/fieldpath"
)

// ResourceConfig represents instructions to grm-generate on how to deal with a
// particular resource.
type ResourceConfig struct {
	// Fields contains a map, keyed by field path, of field configurations
	Fields map[string]*FieldConfig `json:"fields"`
	// AWS returns the AWS-specific resource configuration
	AWS *AWSResourceConfig `json:"aws,omitempty"`
}

// GetFieldConfigs returns a map, keyed by field path, of field configurations
func (c *ResourceConfig) GetFieldConfigs() map[string]*FieldConfig {
	if c == nil || len(c.Fields) == 0 {
		return map[string]*FieldConfig{}
	}
	return c.Fields
}

// GetFieldConfig returns the FieldConfig for a specified field path. This
// method uses case-insensitive matching AND takes into account any renames
// that a field might have. If the supplied path matched a renamed field,
// returns the *renamed* field path as the second return value. If the field is
// not renamed, nil is returned for the second return value.
//
// For example, assume the following configuration snippet:
//
// ```yaml
// resources:
//
//	Bucket:
//	  fields:
//	    Name:
//	      renames:
//	        - Bucket
//
// ```
//
// Calling Bucket ResourceConfig's GetFieldConfig("Bucket") would return
// the FieldConfig struct for the "Name" field, since it has renames for
// "Bucket" along with a fieldpath.FromString("Name")
func (c *ResourceConfig) GetFieldConfig(
	path *fieldpath.Path,
) (*FieldConfig, *fieldpath.Path) {
	if c == nil || len(c.Fields) == 0 {
		return nil, nil
	}
	// First try a simple match on the whole stringified path...
	pathString := path.String()
	for searchPath, fc := range c.Fields {
		if strings.EqualFold(pathString, searchPath) {
			return fc, nil
		}
	}
	// Now check to see if there are any renames for each part of the supplied
	// path
	front := path.Front()
	for pathStr, fc := range c.Fields {
		for _, rename := range fc.Renames {
			if strings.EqualFold(front, rename) {
				return fc, fieldpath.FromString(pathStr)
			}
		}
	}
	return nil, nil
}

// ForAWS returns the AWS-specific resource configuration
func (c *ResourceConfig) ForAWS() *AWSResourceConfig {
	if c != nil && c.AWS != nil {
		return c.AWS
	}
	return nil
}

// AWSResourceConfig contains AWS-specific configuration options for this
// resource
type AWSResourceConfig struct {
	// Operations contains a list of overrides for this resource's operations
	Operations []*AWSResourceOperationConfig `json:"operations"`
}

// AWSResourceOperationConfig instructs the generator which AWS SDK Operation
// to use for which type of operation for this resource.
type AWSResourceOperationConfig struct {
	// Type contains the stringified OpType, e.g. "create" or "READ_ONE"
	Type string `json:"type"`
	// ID contains the ID/name of the AWS SDK Operation that will serve as the
	// OpType for this resource.
	ID string `json:"id"`
}
