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

package model

import (
	"github.com/gertd/go-pluralize"
)

// Kind describes a provider-specific, service-specific type of Resource
type Kind struct {
	// CloudProvider contains the short name of the cloud provider exposing
	// this type of Resource
	CloudProvider string
	// ServiceName contains the short name of the service exposing this type of
	// Resource
	Service string
	// Name contains the camel-cased name of the resource (i.e. the Kind, in
	// Kubernetes speak).
	//
	// Note that the combination of CloudProvider, Service and Name is a unique
	// identifier for this type of Resource.
	Name string
	// PluralName contains the camel-cased name of the pluralized resource.
	//
	// Note that the combination of CloudProvider, Service and PluralName is a
	// unique identifier for this type of Resource.
	PluralName string
}

// NewKind returns a new Kind that describes the type of a single top-level
// resource in a cloud service API
func NewKind(
	cloudProvider string,
	service string,
	name string,
) Kind {
	pluralize := pluralize.NewClient()
	pluralName := pluralize.Plural(name)
	return Kind{
		CloudProvider: cloudProvider,
		Service:       service,
		Name:          name,
		PluralName:    pluralName,
	}
}
