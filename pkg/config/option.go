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

package config

import "strings"

type option struct {
	path string
	yaml string
}

func WithPath(path string) option {
	return option{
		path: path,
	}
}

func WithYAML(yaml string) option {
	// Let's enable a better UX by automatically stripping leading and trailing
	// newlines and converting TAB characters to two spaces.
	corrected := strings.TrimSpace(yaml)
	corrected = strings.ReplaceAll(corrected, "\t", "    ")
	return option{
		yaml: corrected,
	}
}

// mergeOptions merges any supplied option values with any defaults and returns
// a single option
func mergeOptions(opts []option) option {
	res := option{}
	for _, opt := range opts {
		if opt.path != "" {
			res.path = opt.path
		}
		if opt.yaml != "" {
			res.yaml = opt.yaml
		}
	}
	// now process the defaults...
	return res
}
