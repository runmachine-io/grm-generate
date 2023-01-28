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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/anydotcloud/grm-generate/pkg/discover"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

// discoverer is a helper struct that helps work with the aws-sdk-go models and
// API model loader. It implements the `pkg/discover.DiscoversResources`
// interface.
type discoverer struct {
	cfg      ackgenconfig.Config
	basePath string
	loader   *awssdkmodel.Loader
}

func (d *discoverer) DiscoverResources(
	ctx context.Context,
) []model.ResourceDefinition {
	res := model.ResourceDefinition{}
	return res
}

// New returns a new DiscoversResources implementer for AWS resources
func New(
	basePath string,
	cfg ackgenconfig.Config,
) discover.DiscoversResources {
	return &discoverer{
		cfg:      cfg,
		basePath: basePath,
		loader: &awssdkmodel.Loader{
			BaseImport:            basePath,
			IgnoreUnsupportedAPIs: true,
		},
	}
}
