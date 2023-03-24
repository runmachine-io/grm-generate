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
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
)

// FieldDefinition represents a type of Field in a Resource's Schema. Note that
// multiple Fields can have the same FieldDefinition but will never have the
// same FieldDefinition *and* Path.
type FieldDefinition struct {
	// Type is the underlying type of the field.
	Type schema.FieldType `json:"type"`
	// ElementType is the type of the list's elements.
	//
	// If Type is FieldTypeList, the ElementType() method is guaranteed to
	// return the type of the list element. If Type is not FieldTypeList,
	// ElementType is guaranteed to be FieldTypeNil.
	ElementType schema.FieldType `json:"element_type,omitempty"`
	// ValueType is the type of the map's values.
	//
	// If Type is FieldTypeMap, the ValueType() method is guaranteed to return
	// the type of the map values. If Type is not FieldTypeMap, ValueType will
	// always return FieldTypeNil
	ValueType schema.FieldType `json:"value_type,omitempty"`
	// KeyType is the type of the map's keys.
	//
	// If Type is FieldTypeMap, the KeyType() method is guaranteed to return
	// the type of the map keys. If Type is not FieldTypeMap, KeyType will
	// always return FieldTypeNil
	KeyType schema.FieldType `json:"key_type,omitempty"`
	// MemberFieldDefinitions is a map, keyed by member field name, of nested
	// FieldDefinitions when this Field has a Type of FieldTypeStruct. Returns
	// nil when Type is not FieldTypeStruct.
	MemberFieldDefinitions map[string]*FieldDefinition `json:"member_field_definitions,omitempty"`
	// IsRequired is true if the field is required to be set by the user
	IsRequired bool `json:"is_required,omitempty"`
	// IsReadOnly is true if the field is not settable by the user
	IsReadOnly bool `json:"is_read_only,omitempty"`
	// IsImmutable is true if the field cannot be changed once set
	IsImmutable bool `json:"is_immutable,omitempty"`
	// IsLateInitialized is true if the field is "late initialized"
	// with a service-side default value
	IsLateInitialized bool `json:"is_late_initialized,omitempty"`
	// IsSecret is true if the field contains secret information
	IsSecret bool `json:"is_secret,omitempty"`
	// References contains the Kind for a referred type if the field contains a
	// reference to another resource, or nil otherwise.
	//
	// For example, consider a Resource `rds.aws/DBInstance` with a field
	// `Subnets`. This field contains EC2 VPC Subnet identifiers. The Type() of
	// this field would be FieldTypeList. The ElementType() of this field would
	// be FieldTypeString. The References() of this field would return a Kind
	// containing "ec2.aws/Subnet".
	References *Kind `json:"references,omitempty"`
}
