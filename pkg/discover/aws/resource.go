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
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/anydotcloud/grm-generate/pkg/config"
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
		// Finally, we "flatten" the nested field definitions into a
		// single-dimension map of Fields in the ResourceDefinition
		for pathString, f := range r.Fields {
			path := fieldpath.FromString(pathString)
			if path.Size() == 1 {
				flattenField(ctx, rc, r, f, path)
			}
		}
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
		fd := getFieldDefinition(ctx, path, cfg, nil)
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

// flattenField recurses through the supplied field definition's member fields
// ensuring that the resource definition's Fields map contains all nested
// struct fields, keyed by field path.
func flattenField(
	ctx context.Context,
	cfg *config.ResourceConfig,
	r *model.ResourceDefinition,
	f *model.Field,
	path *fieldpath.Path,
) {
	for memberName, memberDef := range f.Definition.MemberFieldDefinitions {
		memberPath := path.Copy()
		memberPath.PushBack(memberName)
		memberField := r.GetField(memberPath)
		if memberField == nil {
			fieldNames := names.New(memberName)
			fc := cfg.GetFieldConfig(memberPath)
			memberField = model.NewField(fieldNames, memberPath, fc, memberDef)
			r.Fields[memberPath.String()] = memberField
		}
		flattenField(ctx, cfg, r, memberField, memberPath)
	}
}
