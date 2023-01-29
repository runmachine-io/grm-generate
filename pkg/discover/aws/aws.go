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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/anydotcloud/grm-generate/pkg/discover"
	"github.com/anydotcloud/grm-generate/pkg/git"
	"github.com/anydotcloud/grm-generate/pkg/log"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

// discoverer is a helper struct that helps work with the aws-sdk-go models and
// API model loader. It implements the `pkg/discover.DiscoversResources`
// interface.
type discoverer struct {
	opts   option
	loader *awssdkmodel.Loader
	repo   *git.Repository
}

func (d *discoverer) DiscoverResources(
	ctx context.Context,
) ([]model.ResourceDefinition, error) {
	var err error
	l := log.FromContext(ctx)
	if d.repo == nil {
		l.Debug("loading git repository", "cache_path", d.opts.cachePath)
		d.repo, err = git.Open(d.opts.cachePath)
		if err != nil {
			return nil, fmt.Errorf(
				"error loading repository from %s: %v",
				d.opts.cachePath, err,
			)
		}
	}
	d.loader = &awssdkmodel.Loader{
		BaseImport:            d.opts.cachePath,
		IgnoreUnsupportedAPIs: true,
	}
	res := []model.ResourceDefinition{}
	return res, nil
}

// New returns a new DiscoversResources implementer for AWS resources
func New(
	opts ...option,
) discover.DiscoversResources {
	return &discoverer{
		opts: mergeOptions(opts),
	}
}
