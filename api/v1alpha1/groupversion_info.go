// Copyright 2024 Stefan Prodan.
// SPDX-License-Identifier: AGPL-3.0

// Package v1alpha1 contains API Schema definitions for the fluxcd v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=fluxcd.controlplane.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "fluxcd.controlplane.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
