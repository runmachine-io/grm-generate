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

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anydotcloud/grm-generate/pkg/config"
)

var (
	emptyConfig  = config.New()
	bucketConfig = config.New(
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
)

func TestGetResourceConfigs(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		cfg  *config.Config
		exp  map[string]*config.ResourceConfig
	}{
		{
			"Nil config returns empty map",
			nil,
			map[string]*config.ResourceConfig{},
		},
		{
			"Empty config returns empty map",
			emptyConfig,
			map[string]*config.ResourceConfig{},
		},
		{
			"Bucket config returns map with single key",
			bucketConfig,
			map[string]*config.ResourceConfig{
				"Bucket": &config.ResourceConfig{
					Fields: map[string]*config.FieldConfig{
						"Name": &config.FieldConfig{
							Renames: []string{
								"Bucket",
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		assert.Equal(test.exp, test.cfg.GetResourceConfigs())
	}
}
