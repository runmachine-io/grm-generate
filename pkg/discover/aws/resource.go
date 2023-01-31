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
	"github.com/anydotcloud/grm-generate/pkg/log"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

// getResourceDefinitionsForService returns a slice of `ResourceDefinition`
// structs that describe the top-level resources discovered for a supplied AWS
// service API
func getResourceDefinitionsForService(
	ctx context.Context,
	service string, // the service package name
	api *awssdkmodel.API,
	cfg *config.Config,
) ([]*model.ResourceDefinition, error) {
	res := []*model.ResourceDefinition{}

	resOpMap := getResourceOperationMap(ctx, api, cfg)

	for resName, ops := range resOpMap {
		// For now, only care about resources with CREATE operations...
		if _, found := ops[OpTypeCreate]; !found {
			continue
		}
		resNames := names.New(resName)
		kind := model.NewKind("aws", service, resNames.Camel)
		rc := cfg.GetResourceConfig(resName)
		fields, err := getFieldsForResource(ctx, rc, resName, ops)
		if err != nil {
			return nil, err
		}
		r := model.NewResourceDefinition(rc, kind, fields)
		res = append(res, r)
	}
	return res, nil
}

// getFieldsForResource returns a map, keyed by field path, of Field objects
// that describe the supplied resource's fields. Fields are collected by
// looking at the supplied FieldConfig structs and examining the set of
// Operations involving the resource.
func getFieldsForResource(
	ctx context.Context,
	cfg *config.ResourceConfig,
	resName string,
	ops map[OpType]*awssdkmodel.Operation,
) (map[string]*model.Field, error) {
	res := map[string]*model.Field{}

	// We first iterate over any FieldConfigs listed in our ResourceConfig. The
	// FieldConfig structs will tell us whether a field has been renamed from
	// the original AWS API shape.
	for pathString, fc := range cfg.GetFieldConfigs() {
		path := fieldpath.FromString(pathString)
		fieldName := path.Back()
		fieldNames := names.New(fieldName)
		fd := getFieldDefinition(ctx, path, fc, nil)
		f := model.NewField(fieldNames, path, fc, fd)
		res[pathString] = f
	}

	// We start with the Create operation's input and output shape. Members of
	// the input shape are user-settable. Members of the output shape that are
	// not in the input shape are read-only.
	if createOp, found := ops[OpTypeCreate]; found {
		inputShape := createOp.InputRef.Shape
		if inputShape == nil {
			msg := fmt.Sprintf(
				"expected non-nil Input shape for createOp %s.",
				createOp.Name,
			)
			panic(msg)
		}

		for memberName, memberShapeRef := range inputShape.MemberRefs {
			if memberShapeRef.Shape == nil {
				msg := fmt.Sprintf(
					"expected non-nil Shape for member %s in inputShape %s.",
					inputShape.ShapeName, memberName,
				)
				panic(msg)
			}
			path := fieldpath.FromString(memberName)
			// NOTE(jaypipes): ResourceConfig.GetFieldConfig accounts for
			// renamed fields...
			fc := cfg.GetFieldConfig(path)
			if fc != nil {
				// The field is already discovered...
				continue
			}
			fieldNames := names.New(memberName)
			fd := getFieldDefinition(ctx, path, nil, memberShapeRef)
			f := model.NewField(fieldNames, path, fc, fd)
			res[fieldNames.Camel] = f
		}
	}
	return res, nil
}

// getFieldDefinition collects information on the field's definition by
// examining both the FieldConfig and the AWS SDK model ShapeRef.
func getFieldDefinition(
	ctx context.Context,
	path *fieldpath.Path,
	cfg *config.FieldConfig,
	shapeRef *awssdkmodel.ShapeRef,
) *model.FieldDefinition {
	l := log.FromContext(ctx)
	def := &model.FieldDefinition{
		Type:        schema.FieldTypeUnknown,
		ValueType:   schema.FieldTypeNil,
		KeyType:     schema.FieldTypeNil,
		ElementType: schema.FieldTypeNil,
	}
	// First try to determine any type information from the field config
	if cfg != nil {
		if cfg.IsReadOnly != nil {
			def.IsReadOnly = *cfg.IsReadOnly
		}
		if cfg.IsRequired != nil {
			def.IsRequired = *cfg.IsRequired
		}
		if cfg.IsImmutable != nil {
			def.IsImmutable = *cfg.IsImmutable
		}
		if cfg.IsSecret != nil {
			def.IsSecret = *cfg.IsSecret
		}
		if cfg.Type != nil {
			def.Type = schema.StringToFieldType(*cfg.Type)
		}
		if cfg.ElementType != nil {
			def.ElementType = schema.StringToFieldType(*cfg.ElementType)
		}
		if cfg.KeyType != nil {
			def.KeyType = schema.StringToFieldType(*cfg.KeyType)
		}
		if cfg.ValueType != nil {
			def.ValueType = schema.StringToFieldType(*cfg.ValueType)
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
		// this is a pointer to the "parent" containing Shape when the field being
		// processed here is a structure or a list/map of structures.
		// var containerShape *awssdkmodel.Shape = shape
		switch shape.Type {
		case "structure":
			l.Debug("skipping struct field", "path", path.String())
		case "list":
			l.Debug("skipping list field", "path", path.String())
		case "map":
			l.Debug("skipping map field", "path", path.String())
		case "timestamp":
			def.Type = schema.FieldTypeTime
		case "string", "character":
			def.Type = schema.FieldTypeString
		case "boolean":
			def.Type = schema.FieldTypeBool
		case "byte", "short", "integer", "long":
			def.Type = schema.FieldTypeInt
		}
	}
	return def
}
