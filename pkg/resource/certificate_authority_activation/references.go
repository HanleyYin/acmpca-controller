// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package certificate_authority_activation

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"

	svcapitypes "github.com/aws-controllers-k8s/acmpca-controller/apis/v1alpha1"
)

// ClearResolvedReferences removes any reference values that were made
// concrete in the spec. It returns a copy of the input AWSResource which
// contains the original *Ref values, but none of their respective concrete
// values.
func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) acktypes.AWSResource {
	ko := rm.concreteResource(res).ko.DeepCopy()

	if ko.Spec.CertificateAuthorityRef != nil {
		ko.Spec.CertificateAuthorityARN = nil
	}

	return &resource{ko}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve those
// reference field(s) into their respective target field(s). It returns a
// copy of the input AWSResource with resolved reference(s), a boolean which
// is set to true if the resource contains any references (regardless of if
// they are resolved successfully) and an error if the passed AWSResource's
// reference field(s) could not be resolved.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, bool, error) {
	ko := rm.concreteResource(res).ko

	resourceHasReferences := false
	err := validateReferenceFields(ko)
	if fieldHasReferences, err := rm.resolveReferenceForCertificateAuthorityARN(ctx, apiReader, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}

	return &resource{ko}, resourceHasReferences, err
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.CertificateAuthorityActivation) error {

	if ko.Spec.CertificateAuthorityRef != nil && ko.Spec.CertificateAuthorityARN != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("CertificateAuthorityARN", "CertificateAuthorityRef")
	}
	if ko.Spec.CertificateAuthorityRef == nil && ko.Spec.CertificateAuthorityARN == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("CertificateAuthorityARN", "CertificateAuthorityRef")
	}
	return nil
}

// resolveReferenceForCertificateAuthorityARN reads the resource referenced
// from CertificateAuthorityRef field and sets the CertificateAuthorityARN
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error
func (rm *resourceManager) resolveReferenceForCertificateAuthorityARN(
	ctx context.Context,
	apiReader client.Reader,
	ko *svcapitypes.CertificateAuthorityActivation,
) (hasReferences bool, err error) {
	if ko.Spec.CertificateAuthorityRef != nil && ko.Spec.CertificateAuthorityRef.From != nil {
		hasReferences = true
		arr := ko.Spec.CertificateAuthorityRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: CertificateAuthorityRef")
		}
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
		}
		obj := &svcapitypes.CertificateAuthority{}
		if err := getReferencedResourceState_CertificateAuthority(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.CertificateAuthorityARN = (*string)(obj.Status.ACKResourceMetadata.ARN)
	}

	return hasReferences, nil
}

// getReferencedResourceState_CertificateAuthority looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_CertificateAuthority(
	ctx context.Context,
	apiReader client.Reader,
	obj *svcapitypes.CertificateAuthority,
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"CertificateAuthority",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"CertificateAuthority",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"CertificateAuthority",
			namespace, name)
	}
	if obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"CertificateAuthority",
			namespace, name,
			"Status.ACKResourceMetadata.ARN")
	}
	return nil
}
