{{- template "boilerplate" }}

package schema

import (
	"github.com/anydotcloud/grm/pkg/types/resource/schema"
)

type kind struct {}

// Service returns the name of the cloud service this resource is
// associated with.
//
// For AWS resources, the string returned matches the service package name
// in aws-sdk-go.
func (k *kind) Service() string {
    return "{{ .Service }}"
}

// Name returns the camel-cased name of the resource (i.e. the Kind, in Kubernetes
// speak).
//
// Note that the combination of Service and Name is a unique identifier for
// this type of Resource.
func (k *kind) Name() string {
    return "{{ .Name }}"
}

// PluralName returns camel-cased name of the pluralized resource.
//
// Note that the combination of Service and PluralName is a unique identifier for
// this type of Resource.
func (k *kind) PluralName() string {
    return "{{ .PluralName }}"
}

var Kind schema.Kind = &kind{}
