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
	"io/ioutil"
	"os"
	"path/filepath"

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
	opts option
	repo *git.Repository
	// apis is a map, keyed by service model package name, of API structs
	// representing the operations and shapes of that service's API.
	apis map[string]*awssdkmodel.API
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
	if err = d.loadAPIs(ctx); err != nil {
		return nil, err
	}
	res := []model.ResourceDefinition{}
	return res, nil
}

// loadModels loads API structs for each service package for which we we are
// discovering resources.
func (d *discoverer) loadAPIs(
	ctx context.Context,
) error {
	if len(d.opts.services) == 0 {
		return nil
	}
	l := log.FromContext(ctx)
	loader := &awssdkmodel.Loader{
		BaseImport:            d.opts.cachePath,
		IgnoreUnsupportedAPIs: true,
	}
	modelPaths, err := d.getModelPaths(ctx)
	if err != nil {
		return err
	}
	apis, err := loader.Load(modelPaths)
	if err != nil {
		return err
	}
	// apis is a map, keyed by the service alias, of pointers to aws-sdk-go
	// model API objects
	for service, api := range apis {
		l.Debug("loading API model", "service", service)
		// If we don't do this, we can end up with panic()'s like this:
		// panic: assignment to entry in nil map
		// when trying to execute Shape.GoType().
		//
		// Calling API.ServicePackageDoc() ends up resetting the API.imports
		// unexported map variable...
		_ = api.ServicePackageDoc()
		//sdkapi := model.NewSDKAPI(api, h.APIGroupSuffix)

		//return sdkapi, nil
	}
	return nil
}

// getModelPaths returns a slice of paths to API model definitions for each
// service for which we are discovering resources. The resulting slice is
// sorted.
func (d *discoverer) getModelPaths(
	ctx context.Context,
) ([]string, error) {
	modelAPIsPath := filepath.Join(d.opts.cachePath, "models", "apis")
	fi, err := os.Lstat(modelAPIsPath)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", modelAPIsPath)
	}

	l := log.FromContext(ctx)
	res := make([]string, len(d.opts.services))
	for x, service := range d.opts.services {
		serviceAPIPath := filepath.Join(modelAPIsPath, service)
		fi, err := os.Lstat(serviceAPIPath)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", serviceAPIPath)
		}
		versionDirs, err := ioutil.ReadDir(serviceAPIPath)
		if err != nil {
			return nil, err
		}
		var apiVersion string
		var serviceAPIVersionPath string
		for _, f := range versionDirs {
			apiVersion = f.Name()
			serviceAPIVersionPath = filepath.Join(serviceAPIPath, apiVersion)
			fi, err := os.Lstat(serviceAPIVersionPath)
			if err != nil {
				return nil, err
			}
			if !fi.IsDir() {
				return nil, fmt.Errorf("%s is not a directory", serviceAPIVersionPath)
			}
			// We only look at the first version...
			break
		}
		apiModelPath := filepath.Join(serviceAPIPath, apiVersion, "api-2.json")
		fi, err = os.Lstat(apiModelPath)
		if err != nil {
			return nil, err
		}
		if !fi.Mode().IsRegular() {
			return nil, fmt.Errorf("%s is not a regular file", apiModelPath)
		}
		res[x] = apiModelPath
		l.Debug("found API model file", "service", service, "path", apiModelPath)
	}
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
