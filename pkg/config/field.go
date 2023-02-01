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

// FieldConfig represents instructions to grm-generate on how to deal with a
// particular resource field.
type FieldConfig struct {
	// Renames instructs the code generator to consider the field to be a
	// rename of one or more names.
	//
	// For example, suppose we are writing a configuration block for the S3
	// Bucket resource. The CreateBucket's Input shape has a Bucket member and
	// through the normal course of API discovery/inference, the Bucket
	// resource would get a field called "Bucket" added to it. If we wanted to
	// rename that to just "Name", we could do the following:
	//
	// ```yaml
	// resources:
	//   Bucket:
	//     fields:
	//       Name:
	//         renames:
	//           - Bucket
	// ```
	//
	// Any time the generator sees the name "Bucket", it will automatically
	// know that the "Name" field is what should be referred to.
	Renames []string `json:"renames,omitempty"`
	// Type *overrides* the type of the field. This is required for custom
	// fields that are not inferred either as a Create Input/Output shape or
	// via the SourceFieldConfig attribute.
	//
	// As an example, assume you have a Role resource where you want to add a
	// custom field called Policies that is a slice of string pointers.
	//
	// The config snippet might look like this:
	//
	// ```yaml
	// resources:
	//   Role:
	//     fields:
	//       Policies:
	//         type: list
	//         element_type: string
	// ```
	Type *string `json:"type,omitempty"`
	// ElementType *overrides* the element type of the field when the field is
	// of type FieldTypeList.
	ElementType *string `json:"element_type,omitempty"`
	// KeyType *overrides* the key type of the field when the field is
	// of type FieldTypeMap.
	KeyType *string `json:"key_type,omitempty"`
	// ValueType *overrides* the value type of the field when the field is
	// of type FieldTypeMap.
	ValueType *string `json:"value_type,omitempty"`
	// IsReadOnly indicates the field's value can not be set by a user
	IsReadOnly *bool `json:"is_read_only,omitempty"`
	// Required indicates whether this field is a required member or not.
	IsRequired *bool `json:"is_required,omitempty"`
	// IsSecret instructs the code generator that this field's value should be
	// considered a secret
	IsSecret *bool `json:"is_secret,omitempty"`
	// IsImmutable instructs the code generator to treat the field as immutable
	// after resource is initially created.
	IsImmutable *bool `json:"is_immutable,omitempty"`
	// AWS returns the AWS-specific field configuration
	AWS *AWSFieldConfig `json:"aws,omitempty"`
}

// ForAWS returns the AWS-specific field configuration
func (c *FieldConfig) ForAWS() *AWSFieldConfig {
	if c != nil && c.AWS != nil {
		return c.AWS
	}
	return nil
}

// AWSFieldConfig contains AWS-specific configuration options for this
// resource
type AWSFieldConfig struct {
}
