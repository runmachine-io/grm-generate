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
	"github.com/samber/lo"

	"github.com/anydotcloud/grm-generate/pkg/config"
)

const (
	DefaultCachePath = "~/.cache/grm-generate"
)

type option struct {
	cfg           *config.Config
	cachePath     string
	services      []string
	apiModelPaths []string
}

// WithConfig uses the supplied Config as instructions to the discovery code
func WithConfig(cfg *config.Config) option {
	return option{
		cfg: cfg,
	}
}

// WithCachePath uses the supplied cache path for discovery of API models
func WithCachePath(path string) option {
	return option{
		cachePath: path,
	}
}

// WithServices instructs the discovery code which AWS services to discover
func WithServices(services ...string) option {
	return option{
		services: services,
	}
}

// WithAPIModelPaths instructs the discovery code where to find API models.
// Expects absolute filepaths. Overrides the `WithCachePath` option.
func WithAPIModelPaths(apiModelPaths ...string) option {
	return option{
		apiModelPaths: apiModelPaths,
	}
}

// mergeOptions merges any supplied option values with any defaults and returns
// a single option
func mergeOptions(opts []option) option {
	res := option{}
	for _, opt := range opts {
		if opt.cfg != nil {
			res.cfg = opt.cfg
		}
		if opt.cachePath != "" {
			res.cachePath = opt.cachePath
		}
		if len(opt.services) > 0 {
			res.services = lo.Uniq(lo.Union(res.services, opt.services))
		}
		if len(opt.apiModelPaths) > 0 {
			res.apiModelPaths = lo.Uniq(
				lo.Union(res.apiModelPaths, opt.apiModelPaths),
			)
		}
	}
	// now process the defaults...
	if res.cachePath == "" {
		res.cachePath = DefaultCachePath
	}
	return res
}
