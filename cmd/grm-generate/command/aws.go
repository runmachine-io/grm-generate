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
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	discover "github.com/anydotcloud/grm-generate/pkg/discover/aws"
	"github.com/anydotcloud/grm-generate/pkg/model"
)

const (
	awsSDKRepoURL = "https://github.com/aws/aws-sdk-go"
)

// discoverAWSCmd is the command that discovers AWS resource models
var discoverAWSCmd = &cobra.Command{
	Use:   "aws <service>",
	Short: "Discover resource models for an AWS service API",
	RunE:  discoverAWS,
}

func init() {
	discoverCmd.AddCommand(discoverAWSCmd)
}

// discoverAWS reads AWS API definitions and discovers resource models
func discoverAWS(
	cmd *cobra.Command,
	args []string,
) error {
	if len(args) != 1 {
		return fmt.Errorf("please specify the service alias for the AWS service API to generate")
	}
	svcAlias := strings.ToLower(args[0])
	ctx, cancel := newContext(context.Background())
	defer cancel()

	sdkRepoTag := ""
	err := cacheRepo(ctx, optCachePath, awsSDKRepoURL, sdkRepoTag)
	if err != nil {
		return err
	}
	sdkCachePath := filepath.Join(optCachePath, "aws-sdk-go")
	disco := discover.New(
		discover.WithCachePath(sdkCachePath),
		discover.WithServices(svcAlias),
	)
	resources, err := disco.DiscoverResources(ctx)
	if err != nil {
		return err
	}
	switch optOutput {
	case "yaml":
		return printResourceDefinitionsYAML(os.Stdout, resources)
	case "table":
		return printResourceDefinitionsTable(os.Stdout, resources)
	}
	return nil
}

func printResourceDefinitionsYAML(
	w io.Writer,
	resources []*model.ResourceDefinition,
) error {
	r := struct {
		Resources []*model.ResourceDefinition
	}{resources}
	y, err := yaml.Marshal(&r)
	if err != nil {
		return err
	}
	_, err = w.Write(y)
	return err
}

func printResourceDefinitionsTable(
	w io.Writer,
	resources []*model.ResourceDefinition,
) error {
	table := tablewriter.NewWriter(w)
	headers := []string{
		"Resource",
		"Field",
		"Type",
		"Required?",
	}
	table.SetHeader(headers)
	data := [][]string{}
	for _, r := range resources {
		rname := r.Kind.Name
		for _, path := range r.GetFieldPaths() {
			f := r.GetField(path)
			typ := f.Definition.Type
			data = append(data, []string{
				rname, path.String(), typ.String(), strconv.FormatBool(f.Definition.IsRequired),
			})
		}
	}
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
	return nil
}
