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

	"github.com/anydotcloud/grm/pkg/path/fieldpath"
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/stretchr/testify/assert"

	"github.com/anydotcloud/grm-generate/pkg/config"
	"github.com/anydotcloud/grm-generate/pkg/discover/aws"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

var (
	flatNoTypeConfig = config.New(
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
	flatWithTypeConfig = config.New(
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
	nestedPathConfig = config.New(
		config.WithYAML(`
resources:
  Repository:
    fields:
      Tags.Value:
		type: string
`,
		),
	)
)

func Test_GetFieldDefinition(t *testing.T) {
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
			"field with no type info and nil shape panics",
			"Name",
			flatNoTypeConfig.GetResourceConfig("Bucket"),
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			true,
		},
		{
			"field with string type info and nil shape infers from config",
			"Name",
			flatWithTypeConfig.GetResourceConfig("Bucket"),
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			false,
		},
		{
			"nested field path with string type info and nil shape infers from config",
			"Tags.Value",
			nestedPathConfig.GetResourceConfig("Repository"),
			nil,
			&model.FieldDefinition{
				Type: schema.FieldTypeString,
			},
			false,
		},
		{
			"case-insensitive nested field path matching",
			"tags.value",
			nestedPathConfig.GetResourceConfig("Repository"),
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
					aws.GetFieldDefinition(ctx, path, test.cfg, test.shapeRef)
				},
			)
		} else {
			got := aws.GetFieldDefinition(ctx, path, test.cfg, test.shapeRef)
			assert.Equal(test.exp, got)
		}
	}
}
