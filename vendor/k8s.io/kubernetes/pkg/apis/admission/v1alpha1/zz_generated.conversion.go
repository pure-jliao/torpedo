// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	pkg_admission "k8s.io/apiserver/pkg/admission"
	admission "k8s.io/kubernetes/pkg/apis/admission"
	unsafe "unsafe"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_AdmissionReview_To_admission_AdmissionReview,
		Convert_admission_AdmissionReview_To_v1alpha1_AdmissionReview,
		Convert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec,
		Convert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec,
		Convert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus,
		Convert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus,
	)
}

func autoConvert_v1alpha1_AdmissionReview_To_admission_AdmissionReview(in *AdmissionReview, out *admission.AdmissionReview, s conversion.Scope) error {
	if err := Convert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_AdmissionReview_To_admission_AdmissionReview is an autogenerated conversion function.
func Convert_v1alpha1_AdmissionReview_To_admission_AdmissionReview(in *AdmissionReview, out *admission.AdmissionReview, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdmissionReview_To_admission_AdmissionReview(in, out, s)
}

func autoConvert_admission_AdmissionReview_To_v1alpha1_AdmissionReview(in *admission.AdmissionReview, out *AdmissionReview, s conversion.Scope) error {
	if err := Convert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_admission_AdmissionReview_To_v1alpha1_AdmissionReview is an autogenerated conversion function.
func Convert_admission_AdmissionReview_To_v1alpha1_AdmissionReview(in *admission.AdmissionReview, out *AdmissionReview, s conversion.Scope) error {
	return autoConvert_admission_AdmissionReview_To_v1alpha1_AdmissionReview(in, out, s)
}

func autoConvert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec(in *AdmissionReviewSpec, out *admission.AdmissionReviewSpec, s conversion.Scope) error {
	out.Kind = in.Kind
	if err := runtime.Convert_runtime_RawExtension_To_runtime_Object(&in.Object, &out.Object, s); err != nil {
		return err
	}
	if err := runtime.Convert_runtime_RawExtension_To_runtime_Object(&in.OldObject, &out.OldObject, s); err != nil {
		return err
	}
	out.Operation = pkg_admission.Operation(in.Operation)
	out.Name = in.Name
	out.Namespace = in.Namespace
	out.Resource = in.Resource
	out.SubResource = in.SubResource
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.UserInfo, &out.UserInfo, 0); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec is an autogenerated conversion function.
func Convert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec(in *AdmissionReviewSpec, out *admission.AdmissionReviewSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdmissionReviewSpec_To_admission_AdmissionReviewSpec(in, out, s)
}

func autoConvert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec(in *admission.AdmissionReviewSpec, out *AdmissionReviewSpec, s conversion.Scope) error {
	out.Kind = in.Kind
	out.Name = in.Name
	out.Namespace = in.Namespace
	if err := runtime.Convert_runtime_Object_To_runtime_RawExtension(&in.Object, &out.Object, s); err != nil {
		return err
	}
	if err := runtime.Convert_runtime_Object_To_runtime_RawExtension(&in.OldObject, &out.OldObject, s); err != nil {
		return err
	}
	out.Operation = pkg_admission.Operation(in.Operation)
	out.Resource = in.Resource
	out.SubResource = in.SubResource
	// TODO: Inefficient conversion - can we improve it?
	if err := s.Convert(&in.UserInfo, &out.UserInfo, 0); err != nil {
		return err
	}
	return nil
}

// Convert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec is an autogenerated conversion function.
func Convert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec(in *admission.AdmissionReviewSpec, out *AdmissionReviewSpec, s conversion.Scope) error {
	return autoConvert_admission_AdmissionReviewSpec_To_v1alpha1_AdmissionReviewSpec(in, out, s)
}

func autoConvert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus(in *AdmissionReviewStatus, out *admission.AdmissionReviewStatus, s conversion.Scope) error {
	out.Allowed = in.Allowed
	out.Result = (*v1.Status)(unsafe.Pointer(in.Result))
	return nil
}

// Convert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus is an autogenerated conversion function.
func Convert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus(in *AdmissionReviewStatus, out *admission.AdmissionReviewStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_AdmissionReviewStatus_To_admission_AdmissionReviewStatus(in, out, s)
}

func autoConvert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus(in *admission.AdmissionReviewStatus, out *AdmissionReviewStatus, s conversion.Scope) error {
	out.Allowed = in.Allowed
	out.Result = (*v1.Status)(unsafe.Pointer(in.Result))
	return nil
}

// Convert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus is an autogenerated conversion function.
func Convert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus(in *admission.AdmissionReviewStatus, out *AdmissionReviewStatus, s conversion.Scope) error {
	return autoConvert_admission_AdmissionReviewStatus_To_v1alpha1_AdmissionReviewStatus(in, out, s)
}
