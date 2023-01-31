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

// ResourceConfig represents instructions to grm-generate on how to deal with a
// particular resource.
type ResourceConfig struct {
	// AWS returns the AWS-specific resource configuration
	AWS *AWSResourceConfig `json:"aws,omitempty"`
}

// ForAWS returns the AWS-specific resource configuration
func (c *ResourceConfig) ForAWS() *AWSResourceConfig {
	if c != nil && c.AWS != nil {
		return c.AWS
	}
	return nil
}

// AWSResourceConfig contains AWS-specific configuration options for this
// resource
type AWSResourceConfig struct {
	// Operations contains a map containing overrides for this resource's
	// operations
	Operations []*AWSResourceOperationConfig `json:"operations"`
}

// AWSResourceOperationConfig instructs the generator which AWS SDK Operation
// to use for which type of operation for this resource.
type AWSResourceOperationConfig struct {
	// Type contains the stringified OpType, e.g. "create" or "READ_ONE"
	Type string `json:"type"`
	// ID contains the ID/name of the AWS SDK Operation that will serve as the
	// OpType for this resource.
	ID string `json:"id"`
}
