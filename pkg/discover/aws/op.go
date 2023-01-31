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
	"fmt"
	"strings"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"

	"github.com/anydotcloud/grm-generate/pkg/config"
)

type OpType int

const (
	OpTypeUnknown OpType = iota
	OpTypeCreate
	OpTypeCreateBatch
	OpTypeDelete
	OpTypeReplace
	OpTypeUpdate
	OpTypeAddChild
	OpTypeAddChildren
	OpTypeRemoveChild
	OpTypeRemoveChildren
	OpTypeGet
	OpTypeList
	OpTypeGetAttributes
	OpTypeSetAttributes
)

type resourceOperationMap map[string]map[OpType]*awssdkmodel.Operation

// GetOperationsForResource returns a map, keyed by OpType, for a supplied
// resource. Resource name matching is case-insensitive.
func (m resourceOperationMap) GetOperationsForResource(
	resName string,
) *map[OpType]*awssdkmodel.Operation {
	for name, opMap := range m {
		if strings.EqualFold(name, resName) {
			return &opMap
		}
	}
	return nil
}

// getResourceOperationMap returns a map, keyed by the resource name, of maps,
// keyed by OpType, of aws-sdk-go private/model/api.Operation struct pointers
// that describe that Operation for that resource.
func getResourceOperationMap(
	ctx context.Context,
	api *awssdkmodel.API,
	cfg *config.Config,
) resourceOperationMap {
	// create an index of Operations by resource name and operation type
	res := resourceOperationMap{}
	for opID, op := range api.Operations {
		opType, resName := getOpTypeAndResourceNameFromOpID(opID, cfg)
		resOps := res.GetOperationsForResource(resName)
		if resOps == nil {
			resOps = &map[OpType]*awssdkmodel.Operation{}
		}
		(*resOps)[opType] = op
		res[resName] = *resOps
	}

	// We need to do a second pass over the operation overrides because some
	// APIs have multiple operations of a particular OpType and we need to
	// always be sure that the overridden operation is the one that we use
	// during inference.
	//
	// An example of this is the Kinesis API which has two OpTypeGet
	// Operations: DescribeStream and DescribeStreamSummary. We want the latter
	// only and list that in our `operations:` configuration value.
	for resName, rc := range cfg.GetResourceConfigs() {
		arc := rc.ForAWS()
		if arc == nil {
			continue
		}
		resOps := res.GetOperationsForResource(resName)
		if resOps == nil {
			resOps = &map[OpType]*awssdkmodel.Operation{}
		}
		for x, aroc := range arc.Operations {
			opID := aroc.ID
			opType := getOpTypeFromString(aroc.Type)
			if opType == OpTypeUnknown {
				msg := fmt.Sprintf(
					"operation type %s in config 'resources[%s].aws.operations[%d]:' "+
						"is unknown",
					aroc.Type, resName, x,
				)
				panic(msg)
			}
			op, found := api.Operations[opID]
			if !found {
				msg := fmt.Sprintf(
					"operation %s in config 'resources[%s].aws.operations:' "+
						"does not exist in API model.",
					opID, resName,
				)
				panic(msg)
			}
			(*resOps)[opType] = op
		}
	}
	return res
}

// getOpTypeAndResourceNameFromOpID guesses the resource name and type of
// operation from the OperationID
func getOpTypeAndResourceNameFromOpID(
	opID string,
	cfg *config.Config,
) (OpType, string) {
	pluralize := pluralize.NewClient()
	if strings.HasPrefix(opID, "CreateOrUpdate") {
		return OpTypeReplace, strings.TrimPrefix(opID, "CreateOrUpdate")
	} else if strings.HasPrefix(opID, "BatchCreate") {
		resName := strings.TrimPrefix(opID, "BatchCreate")
		if pluralize.IsPlural(resName) {
			// Do not singularize "pluralized singular" resources
			// like EC2's DhcpOptions, if defined in generator config's list of
			// resources.
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeCreateBatch, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreateBatch, resName
	} else if strings.HasPrefix(opID, "CreateBatch") {
		resName := strings.TrimPrefix(opID, "CreateBatch")
		if pluralize.IsPlural(resName) {
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeCreateBatch, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreateBatch, resName
	} else if strings.HasPrefix(opID, "Create") {
		resName := strings.TrimPrefix(opID, "Create")
		if pluralize.IsPlural(resName) {
			// If resName exists in the generator configuration's list of
			// resources, then just return OpTypeCreate and the resource name.
			// This handles "pluralized singular" resource names like EC2's
			// DhcpOptions.
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeCreate, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreate, resName
	} else if strings.HasPrefix(opID, "Modify") {
		return OpTypeUpdate, strings.TrimPrefix(opID, "Modify")
	} else if strings.HasPrefix(opID, "Update") {
		return OpTypeUpdate, strings.TrimPrefix(opID, "Update")
	} else if strings.HasPrefix(opID, "Delete") {
		return OpTypeDelete, strings.TrimPrefix(opID, "Delete")
	} else if strings.HasPrefix(opID, "Describe") {
		resName := strings.TrimPrefix(opID, "Describe")
		if pluralize.IsPlural(resName) {
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeList, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeGet, resName
	} else if strings.HasPrefix(opID, "Get") {
		if strings.HasSuffix(opID, "Attributes") {
			resName := strings.TrimPrefix(opID, "Get")
			resName = strings.TrimSuffix(resName, "Attributes")
			return OpTypeGetAttributes, resName
		}
		resName := strings.TrimPrefix(opID, "Get")
		if pluralize.IsPlural(resName) {
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeGet, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeGet, resName
	} else if strings.HasPrefix(opID, "List") {
		resName := strings.TrimPrefix(opID, "List")
		if pluralize.IsPlural(resName) {
			rc := cfg.GetResourceConfig(resName)
			if rc != nil {
				return OpTypeList, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeList, resName
	} else if strings.HasPrefix(opID, "Set") {
		if strings.HasSuffix(opID, "Attributes") {
			resName := strings.TrimPrefix(opID, "Set")
			resName = strings.TrimSuffix(resName, "Attributes")
			return OpTypeSetAttributes, resName
		}
	}
	return OpTypeUnknown, opID
}

// getOpTypeFromString translates a string literal into the associated OpType
func getOpTypeFromString(s string) OpType {
	switch strings.ToLower(s) {
	case "create":
		return OpTypeCreate
	case "createbatch":
		return OpTypeCreateBatch
	case "delete":
		return OpTypeDelete
	case "replace":
		return OpTypeReplace
	case "update":
		return OpTypeUpdate
	case "addchild":
		return OpTypeAddChild
	case "addchildren":
		return OpTypeAddChildren
	case "removechild":
		return OpTypeRemoveChild
	case "removechildren":
		return OpTypeRemoveChildren
	case "get", "readone", "read_one":
		return OpTypeGet
	case "list", "readmany", "read_many":
		return OpTypeList
	case "getattributes", "get_attributes":
		return OpTypeGetAttributes
	case "setattributes", "set_attributes":
		return OpTypeSetAttributes
	}

	return OpTypeUnknown
}
