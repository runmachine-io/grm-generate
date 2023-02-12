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

// GetResourceDefinitionsForService returns a slice of `ResourceDefinition`
// structs that describe the top-level resources discovered for a supplied AWS
// service API
func GetResourceDefinitionsForService(
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
		rd := model.NewResourceDefinition(rc, kind)
		err := AddFieldsToResourceDefinition(ctx, rd, rc, ops)
		if err != nil {
			return nil, err
		}
		res = append(res, rd)
	}
	return res, nil
}

// AddFieldsToResourceDefinition iterates over API Operations and a supplied
// ResourceConfig and adds Fields to the supplied ResourceDefinition, recursing
// down through any nested fields.
func AddFieldsToResourceDefinition(
	ctx context.Context,
	rd *model.ResourceDefinition,
	cfg *config.ResourceConfig,
	ops map[OpType]*awssdkmodel.Operation,
) error {
	rName := rd.Kind.Name

	// We start with the Create operation's input and output shape. Members of
	// the input shape are user-settable. Members of the output shape that are
	// not in the input shape are read-only.
	if createOp, found := ops[OpTypeCreate]; found {
		inputShape := createOp.InputRef.Shape
		if inputShape == nil {
			msg := fmt.Sprintf(
				"processing resource %s found nil Input shape for createOp %s.",
				rName, createOp.Name,
			)
			panic(msg)
		}

		for memberName, memberShapeRef := range inputShape.MemberRefs {
			if memberShapeRef.Shape == nil {
				msg := fmt.Sprintf(
					"processing resource %s found nil Shape for member %s in inputShape %s.",
					rName, inputShape.ShapeName, memberName,
				)
				panic(msg)
			}
			path := fieldpath.FromString(memberName)
			VisitMemberShape(ctx, rd, path, cfg, memberShapeRef)
		}
	}
	return nil
}
