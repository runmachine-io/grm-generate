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
	"github.com/jaypipes/grm/pkg/types/resource/schema"
)

// FieldDefinition represents a type of Field in a Resource's Schema. Note that
// multiple Fields can have the same FieldDefinition but will never have the
// same FieldDefinition *and* Path.
type FieldDefinition struct {
	// Name is the *renamed, cleaned, camel-cased name* of the field
	// definition. This is akin to the aws-sdk-go private/model/api.Shape.Name
	// attribute
	Name string
	// Type is the underlying type of the field.
	Type schema.FieldType
	// ElementType is the type of the list's elements.
	//
	// If Type is FieldTypeList, the ElementType() method is guaranteed to
	// return the type of the list element. If Type is not FieldTypeList,
	// ElementType is guaranteed to be FieldTypeNil.
	ElementType schema.FieldType
	// ValueType is the type of the map's values.
	//
	// If Type is FieldTypeMap, the ValueType() method is guaranteed to return
	// the type of the map values. If Type is not FieldTypeMap, ValueType will
	// always return FieldTypeNil
	ValueType schema.FieldType
	// KeyType is the type of the map's keys.
	//
	// If Type is FieldTypeMap, the KeyType() method is guaranteed to return
	// the type of the map keys. If Type is not FieldTypeMap, KeyType will
	// always return FieldTypeNil
	KeyType schema.FieldType
	// MemberFields is a map, keyed by member field name, of nested Fields
	// when this Field has a Type of FieldTypeStruct. Returns nil when Type is
	// not FieldTypeStruct.
	MemberFields map[string]Field
	// IsReadOnly is true if the field is not settable by the user
	IsReadOnly bool
	// IsImmutable is true if the field cannot be changed once set
	IsImmutable bool
	// IsLateInitialized is true if the field is "late initialized"
	// with a service-side default value
	IsLateInitialized bool
	// IsSecret is true if the field contains secret information
	IsSecret bool
	// References contains the Kind for a referred type if the field contains a
	// reference to another resource, or nil otherwise.
	//
	// For example, consider a Resource `rds.aws/DBInstance` with a field
	// `Subnets`. This field contains EC2 VPC Subnet identifiers. The Type() of
	// this field would be FieldTypeList. The ElementType() of this field would
	// be FieldTypeString. The References() of this field would return a Kind
	// containing "ec2.aws/Subnet".
	References *Kind
}
