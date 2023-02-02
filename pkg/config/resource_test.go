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

	"github.com/anydotcloud/grm/pkg/path/fieldpath"
	"github.com/stretchr/testify/assert"

	"github.com/anydotcloud/grm-generate/pkg/config"
)

var (
	emptyResourceConfig = &config.ResourceConfig{}
	bucketConfig        = s3Config.GetResourceConfig("Bucket")
)

func TestGetFieldConfigs(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		cfg  *config.ResourceConfig
		exp  map[string]*config.FieldConfig
	}{
		{
			"Nil config returns empty map",
			nil,
			map[string]*config.FieldConfig{},
		},
		{
			"Empty config returns empty map",
			emptyResourceConfig,
			map[string]*config.FieldConfig{},
		},
		{
			"Bucket config returns map with single key",
			bucketConfig,
			map[string]*config.FieldConfig{
				"Name": &config.FieldConfig{
					Renames: []string{
						"Bucket",
					},
				},
			},
		},
	}
	for _, test := range tests {
		assert.Equal(test.exp, test.cfg.GetFieldConfigs())
	}
}

func TestGetFieldConfig(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name      string
		fieldPath string
		cfg       *config.ResourceConfig
		exp       *config.FieldConfig
	}{
		{
			"Nil config returns nil",
			"Name",
			nil,
			nil,
		},
		{
			"Empty config returns nil",
			"Name",
			emptyResourceConfig,
			nil,
		},
		{
			"Name returns FieldConfig",
			"Name",
			bucketConfig,
			&config.FieldConfig{
				Renames: []string{
					"Bucket",
				},
			},
		},
		{
			"lowercase Name returns FieldConfig",
			"name",
			bucketConfig,
			&config.FieldConfig{
				Renames: []string{
					"Bucket",
				},
			},
		},
		{
			"renamed from Bucket returns FieldConfig",
			"Bucket",
			bucketConfig,
			&config.FieldConfig{
				Renames: []string{
					"Bucket",
				},
			},
		},
		{
			"renamed from Bucket returns FieldConfig when lowercase rename",
			"bucket",
			bucketConfig,
			&config.FieldConfig{
				Renames: []string{
					"Bucket",
				},
			},
		},
	}
	for _, test := range tests {
		path := fieldpath.FromString(test.fieldPath)
		assert.Equal(test.exp, test.cfg.GetFieldConfig(path))
	}
}
