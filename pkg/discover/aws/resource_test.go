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
	"testing"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm-generate/pkg/discover/aws"
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

func Test_GetResourceDefinitionForService_ECR(t *testing.T) {
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
	assert.Equal(2, len(rds))
}
