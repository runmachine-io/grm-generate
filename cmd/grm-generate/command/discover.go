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

package command

import (
	"github.com/spf13/cobra"
)

// discoverCmd is the command that discovers resources
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover resource models",
}

func init() {
	discoverCmd.PersistentFlags().StringVarP(
		&optOutput, "output", "o", "table",
		"Output in what format?",
	)
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}
