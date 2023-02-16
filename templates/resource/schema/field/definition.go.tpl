{{- template "boilerplate" }}

package field

import (
	"github.com/anydotcloud/grm/pkg/types/resource"
	"github.com/anydotcloud/grm/pkg/types/resource/schema"

	res "{{ .ResourcePackage }}"
)

var (
    memberFields{{ .Name }} = map[string]resource.Field{
{{- range $memberFieldName, $memberFieldTypeName := .MemberFields }}
        "{{ $memberFieldName }}": {{ $memberFieldTypeName }},
{{ end -}}
    }
)

type def{{ .Name }} struct {}

// Type returns the underlying type of the field.
func (d *def{{ .Name }}) Type() schema.FieldType {
	return schema.{{ .FieldType.EnumString }}
}

// ElementType returns the type of the list's elements.
//
// If Type is FieldTypeList, the ElementType() method is guaranteed to
// return the type of the list element. If Type is not FieldTypeList,
// ElementType is guaranteed to be FieldTypeNil.
func (d *def{{ .Name }}) ElementType() schema.FieldType {
	return schema.{{ .ElementType.EnumString }}
}

// ValueType returns the type of the map's values.
//
// If Type is FieldTypeMap, the ValueType() method is guaranteed to return
// the type of the map values. If Type is not FieldTypeMap, ValueType will
// always return FieldTypeNil
func (d *def{{ .Name }}) ValueType() schema.FieldType {
	return schema.{{ .ValueType.EnumString }}
}

// KeyType returns the type of the map's keys.
//
// If Type is FieldTypeMap, the KeyType() method is guaranteed to return
// the type of the map keys. If Type is not FieldTypeMap, KeyType will
// always return FieldTypeNil
func (d *def{{ .Name }}) KeyType() schema.FieldType {
	return schema.{{ .KeyType.EnumString }}
}

// MemberFields returns a map, keyed by member field name, of nested Fields
// when this Field has a Type of FieldTypeStruct. Returns nil when Type is
// not FieldTypeStruct.
func (d *def{{ .Name }}) MemberFields() map[string]schema.Field {
    return memberFields{{ .Name }}
}

// IsRequired returns true if the field is required to be set by the user
func (d *def{{ .Name }}) IsRequired() bool {
	return {{ printf "%b" .IsRequired }}
}

// IsReadOnly returns true if the field is not settable by the user
func (d *def{{ .Name }}) IsReadOnly() bool {
	return {{ printf "%b" .IsReadOnly }}
}

// IsImmutable returns true if the field cannot be changed once set
func (d *def{{ .Name }}) IsImmutable() bool {
	return {{ printf "%b" .IsImmutable }}
}

// IsLateInitialized returns true if the field is "late initialized"
// with a service-side default value
func (d *def{{ .Name }}) IsLateInitialized() bool {
	return {{ printf "%b" .IsLateInitialized }}
}

// IsSecret returns true if the field contains secret information
func (d *def{{ .Name }}) IsSecret() bool {
	return {{ printf "%b" .IsSecret }}
}

// References returns the Kind for a referred type if the field contains a
// reference to another resource, or nil otherwise.
//
// For example, consider a Resource `rds.aws/DBInstance` with a field
// `Subnets`. This field contains EC2 VPC Subnet identifiers. The Type() of
// this field would be FieldTypeList. The ElementType() of this field would
// be FieldTypeString. The References() of this field would return a Kind
// containing "ec2.aws/Subnet".
func (d *def{{ .Name }}) References() resource.Kind {
    // TODO(jaypipes)
	return nil
}

{{ .Documentation }}
var {{ .Name }} resource.Field = &def{{ .Name }}{}
