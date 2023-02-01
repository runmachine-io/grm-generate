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
	"testing"

	"github.com/anydotcloud/grm/pkg/path/fieldpath"
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/stretchr/testify/assert"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

var (
	bucketNameNoTypeConfig = config.New(
		config.WithYAML(`
resources:
  Bucket:
    fields:
      Name:
        renames:
         - Bucket
`,
		),
	)
	bucketNameWithTypeConfig = config.New(
		config.WithYAML(`
resources:
  Bucket:
    fields:
      Name:
		renames:
         - Bucket
		type: string
`,
		),
	)
)

func Test_getFieldDefinition(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		path     string
		cfg      *config.ResourceConfig
		shapeRef *awssdkmodel.ShapeRef
		exp      *model.FieldDefinition
		expPanic bool
	}{
		{
			"nil config and nil shape panics",
			"FieldName",
			nil,
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			true,
		},
		{
			"config with no type info and nil shape panics",
			"Name",
			bucketNameNoTypeConfig.GetResourceConfig("Bucket"),
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			true,
		},
		{
			"config with string type info and nil shape infers from config",
			"Name",
			bucketNameWithTypeConfig.GetResourceConfig("Bucket"),
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			false,
		},
	}
	ctx := context.TODO()
	for _, test := range tests {
		path := fieldpath.FromString(test.path)
		if test.expPanic {
			assert.Panics(
				func() {
					getFieldDefinition(ctx, path, test.cfg, test.shapeRef)
				},
			)
		} else {
			got := getFieldDefinition(ctx, path, test.cfg, test.shapeRef)
			assert.Equal(test.exp, got)
		}
	}
}
