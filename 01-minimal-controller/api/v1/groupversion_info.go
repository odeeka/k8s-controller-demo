// Package v1 contains API Schema definitions for the learn.example.com v1 API group.
//
// All custom resource types in this step live in this package.
// The package name "v1" matches the API version in the CRD YAML.
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion describes the API group and version for all types in this package.
	// It must match the `group` and version fields in the CRD YAML.
	GroupVersion = schema.GroupVersion{Group: "learn.example.com", Version: "v1"}

	// SchemeBuilder helps register our types with a runtime.Scheme.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme is a convenience function that registers all types in this
	// package into a scheme. We call it in main.go.
	AddToScheme = SchemeBuilder.AddToScheme
)
