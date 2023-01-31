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

	"github.com/anydotcloud/grm/pkg/names"
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
		r := model.NewResourceDefinition(rc, kind)
		res = append(res, r)
	}
	return res, nil
}
