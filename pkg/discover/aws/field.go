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

package aws

import (
	"context"
	"fmt"

	"github.com/anydotcloud/grm/pkg/names"
	"github.com/anydotcloud/grm/pkg/path/fieldpath"
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

// getFieldDefinition collects information on the field's definition by
// examining both the FieldConfig and the AWS SDK model ShapeRef.
func getFieldDefinition(
	ctx context.Context,
	path *fieldpath.Path,
	cfg *config.ResourceConfig,
	shapeRef *awssdkmodel.ShapeRef,
) *model.FieldDefinition {
	def := &model.FieldDefinition{
		Type:        schema.FieldTypeUnknown,
		ValueType:   schema.FieldTypeNil,
		KeyType:     schema.FieldTypeNil,
		ElementType: schema.FieldTypeNil,
	}
	// First try to determine any type information from the field config
	fc := cfg.GetFieldConfig(path)
	if fc != nil {
		if fc.IsReadOnly != nil {
			def.IsReadOnly = *fc.IsReadOnly
		}
		if fc.IsRequired != nil {
			def.IsRequired = *fc.IsRequired
		}
		if fc.IsImmutable != nil {
			def.IsImmutable = *fc.IsImmutable
		}
		if fc.IsSecret != nil {
			def.IsSecret = *fc.IsSecret
		}
		if fc.Type != nil {
			def.Type = schema.StringToFieldType(*fc.Type)
		}
		if fc.ElementType != nil {
			def.ElementType = schema.StringToFieldType(*fc.ElementType)
		}
		if fc.KeyType != nil {
			def.KeyType = schema.StringToFieldType(*fc.KeyType)
		}
		if fc.ValueType != nil {
			def.ValueType = schema.StringToFieldType(*fc.ValueType)
		}
	}

	if def.Type == schema.FieldTypeUnknown {
		if shapeRef == nil {
			msg := fmt.Sprintf(
				"cannot determine field definition/type for %s. "+
					"No field config or shapeRef supplied.",
				path,
			)
			panic(msg)
		}
		// Let's examine the supplied ShapeRef for type information...
		var shape *awssdkmodel.Shape
		if shapeRef != nil {
			shape = shapeRef.Shape
		}
		if shape == nil {
			msg := fmt.Sprintf(
				"expected non-nil Shape for %s shapeRef",
				shapeRef.ShapeName,
			)
			panic(msg)
		}
		def.Type = fieldTypeFromShape(shape)
		// this is a pointer to the "parent" containing Shape when the field being
		// processed here is a structure or a list/map of structures.
		// var containerShape *awssdkmodel.Shape = shape
		switch shape.Type {
		case "list", "map":
			if shape.Type == "list" {
				def.ElementType = fieldTypeFromShape(shape.MemberRef.Shape)
			} else {
				// Currently only map of string keys is supported...
				def.KeyType = schema.FieldTypeString
				def.ValueType = fieldTypeFromShape(shape.ValueRef.Shape)
			}
			// this is a pointer to the "parent" containing Shape when the field being
			// processed here is a structure or a list/map of structures.
			var containerShape *awssdkmodel.Shape = shape

			for {
				// If the field is a slice or map of structs, we want to add
				// MemberFields that describe the list or value struct elements so
				// that a field path can be used to "find" nested struct member
				// fields.
				//
				// For example, the EC2 resource DHCPOptions has a Field called
				// DHCPConfigurations which is of type []*NewDHCPConfiguration
				// where the NewDHCPConfiguration struct contains two fields, Key
				// and Values. If we want to be able to refer to the
				// DHCPOptions.DHCPConfigurations.Values field by field path, we
				// need a Field.MemberField that describes the
				// NewDHCPConfiguration.Values field.
				//
				// Here, we essentially dive down into list or map fields,
				// searching for whether the list or map fields have structure list
				// element or value element types and then rely on the code below
				// to "unpack" those struct member fields.
				if containerShape.Type == "list" {
					containerShape = containerShape.MemberRef.Shape
					continue
				} else if containerShape.Type == "map" {
					containerShape = containerShape.ValueRef.Shape
					continue
				}
				break
			}

			if containerShape.Type == "structure" {
				def.MemberFieldDefinitions = getMemberFieldDefinitions(ctx, cfg, containerShape, path)
			}
		case "structure":
			def.MemberFieldDefinitions = getMemberFieldDefinitions(ctx, cfg, shape, path)
		}
	}
	return def
}

// fieldTypeFromShape returns the schema.FieldType from an aws-sdk-go
// Shape.Type string.
func fieldTypeFromShape(
	s *awssdkmodel.Shape,
) schema.FieldType {
	switch s.Type {
	case "list":
		return schema.FieldTypeList
	case "map":
		return schema.FieldTypeMap
	case "structure":
		return schema.FieldTypeStruct
	case "timestamp":
		return schema.FieldTypeTime
	case "string", "character":
		return schema.FieldTypeString
	case "boolean":
		return schema.FieldTypeBool
	case "byte", "short", "integer", "long":
		return schema.FieldTypeInt
	default:
		return schema.FieldTypeUnknown
	}
}

// getMemberFieldDefinitions returns a map, keyed by normalized field name, of
// a struct field's member field definitions
func getMemberFieldDefinitions(
	ctx context.Context,
	cfg *config.ResourceConfig,
	containerShape *awssdkmodel.Shape, // the "parent" or "containing" shape
	containerPath *fieldpath.Path, // the field path to containing field
) map[string]*model.FieldDefinition {
	defs := map[string]*model.FieldDefinition{}
	for _, memberName := range containerShape.MemberNames() {
		cleanMemberNames := names.New(memberName)
		memberPath := containerPath.Copy()
		memberPath.PushBack(cleanMemberNames.Camel)
		memberShape := containerShape.MemberRefs[memberName]
		memberDef := getFieldDefinition(ctx, memberPath, cfg, memberShape)
		defs[cleanMemberNames.Camel] = memberDef
	}
	return defs
}
