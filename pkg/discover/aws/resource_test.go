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

package aws_test

import (
	"context"
	"sort"
	"strings"
	"testing"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm-generate/pkg/discover/aws"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

func Test_GetResourceDefinitionForService_Panics(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		service string
		api     *awssdkmodel.API
		cfg     *config.Config
	}{
		{
			"nil config and nil API panics",
			"nonexist",
			nil,
			nil,
		},
		{
			"config with no type info and nil API panics",
			"nonexist",
			nil,
			flatNoTypeConfig,
		},
	}
	ctx := context.TODO()
	for _, test := range tests {
		assert.Panics(
			func() {
				aws.GetResourceDefinitionsForService(
					ctx, test.service, test.api, test.cfg,
				)
			},
		)
	}
}

func Test_GetResourceDefinitionForService_Kind(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	ctx := context.TODO()
	service := "ecr"
	api := apis[service]
	require.NotNil(api, "expected non-nil API for ECR service")
	require.Equal(api.PackageName(), "ecr")
	rds, err := aws.GetResourceDefinitionsForService(
		ctx, service, api, nil,
	)
	require.Nil(err)
	numRDs := 2
	assert.Equal(numRDs, len(rds))

	names := make([]string, numRDs)
	pluralNames := make([]string, numRDs)
	cloudProviders := make([]string, numRDs)
	services := make([]string, numRDs)
	for x, rd := range rds {
		names[x] = rd.Kind.Name
		pluralNames[x] = rd.Kind.PluralName
		cloudProviders[x] = rd.Kind.CloudProvider
		services[x] = rd.Kind.Service
	}
	sort.Strings(names)
	expectNames := []string{
		"PullThroughCacheRule",
		"Repository",
	}
	assert.Equal(expectNames, names)

	sort.Strings(pluralNames)
	expectPluralNames := []string{
		"PullThroughCacheRules",
		"Repositories",
	}
	assert.Equal(expectPluralNames, pluralNames)

	assert.Equal([]string{"aws"}, lo.Uniq(cloudProviders))

	assert.Equal([]string{"ecr"}, lo.Uniq(services))
}

func Test_GetResourceDefinitionForService_FieldPaths_NoConfig(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	ctx := context.TODO()
	service := "ecr"
	api := apis[service]
	require.NotNil(api, "expected non-nil API for ECR service")
	require.Equal(api.PackageName(), "ecr")
	rds, err := aws.GetResourceDefinitionsForService(
		ctx, service, api, nil,
	)
	require.Nil(err)

	var repoRD *model.ResourceDefinition
	for _, rd := range rds {
		if strings.EqualFold(rd.Kind.Name, "repository") {
			repoRD = rd
			break
		}
	}
	require.NotNil(repoRD)

	fieldPaths := []string{}
	for _, fPath := range repoRD.GetFieldPaths() {
		fieldPaths = append(fieldPaths, fPath.String())
	}
	sort.Strings(fieldPaths)

	expectFieldPaths := []string{
		"EncryptionConfiguration",
		"EncryptionConfiguration.EncryptionType",
		"EncryptionConfiguration.KMSKey",
		"ImageScanningConfiguration",
		"ImageScanningConfiguration.ScanOnPush",
		"ImageTagMutability",
		"RegistryID",
		"RepositoryName",
		"Tags",
		"Tags.Key",
		"Tags.Value",
	}
	assert.Equal(expectFieldPaths, fieldPaths)
}

func Test_GetResourceDefinitionForService_FieldPaths_RenamingConfig(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	ctx := context.TODO()
	service := "ecr"
	api := apis[service]
	require.NotNil(api, "expected non-nil API for ECR service")
	require.Equal(api.PackageName(), "ecr")

	cfg := config.New(
		config.WithYAML(`
resources:
  Repository:
    fields:
      Name:
		renames:
         - RepositoryName
`,
		),
	)

	rds, err := aws.GetResourceDefinitionsForService(
		ctx, service, api, cfg,
	)
	require.Nil(err)

	var repoRD *model.ResourceDefinition
	for _, rd := range rds {
		if strings.EqualFold(rd.Kind.Name, "repository") {
			repoRD = rd
			break
		}
	}
	require.NotNil(repoRD)

	fieldPaths := []string{}
	for _, fPath := range repoRD.GetFieldPaths() {
		fieldPaths = append(fieldPaths, fPath.String())
	}
	sort.Strings(fieldPaths)

	expectFieldPaths := []string{
		"EncryptionConfiguration",
		"EncryptionConfiguration.EncryptionType",
		"EncryptionConfiguration.KMSKey",
		"ImageScanningConfiguration",
		"ImageScanningConfiguration.ScanOnPush",
		"ImageTagMutability",
		"Name",
		"RegistryID",
		"Tags",
		"Tags.Key",
		"Tags.Value",
	}
	assert.Equal(expectFieldPaths, fieldPaths)
}
