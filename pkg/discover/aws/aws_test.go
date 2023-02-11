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
	"fmt"
	"path/filepath"

	"github.com/anydotcloud/grm-generate/pkg/discover/aws"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

var (
	apiModelDir, _ = filepath.Abs("testdata")
	services       = []string{
		"dynamodb",
		"ec2",
		"ecr",
		"lambda",
		"s3",
	}
	apiModelPaths []string
	apis          map[string]*awssdkmodel.API
)

func init() {
	ctx := context.TODO()
	for _, service := range services {
		apiModelPaths = append(
			apiModelPaths,
			filepath.Join(apiModelDir, fmt.Sprintf("%s-api.json", service)),
		)
	}
	sapis, err := aws.GetAPIs(ctx, apiModelDir, apiModelPaths)
	if err != nil {
		panic(err)
	}
	apis = sapis
}
