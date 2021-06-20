package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func getRecommendedLabels() []string {
	return []string{
		"app.kubernetes.io/name",
		"app.kubernetes.io/instance",
		"app.kubernetes.io/version",
		"app.kubernetes.io/component",
		"app.kubernetes.io/part-of",
		"app.kubernetes.io/managed-by",
		"app.kubernetes.io/created-by",
	}
}

type ValidationHandler struct{}

func (mh *ValidationHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	var pod v1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		glog.Errorf("failed to unmarshal raw pod object: %v", err)
		return admission.Response{
			AdmissionResponse: v1beta1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			},
		}
	}

	glog.Infof("AdmissionReview Kind=%q\n"+
		"Namespace=%q "+
		"Name=%q/%q "+
		"UID=%v "+
		"patchOperation=%v "+
		"UserInfo=%v",
		req.Kind,
		req.Namespace,
		req.Name,
		pod.Name,
		req.UID,
		req.Operation,
		req.UserInfo)
	if isAllowed, msg := hasRecommendedLabels(&pod); !isAllowed {
		return admission.Response{
			AdmissionResponse: v1beta1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Message: msg,
				},
			},
		}
	}
	return admission.Response{
		AdmissionResponse: v1beta1.AdmissionResponse{
			Allowed: true,
		},
	}
}

func hasRecommendedLabels(pod *v1.Pod) (bool, string) {
	recommendedLabels := getRecommendedLabels()
	var msg string
	for _, labelKey := range recommendedLabels {
		_, ok := pod.Labels[labelKey]
		if !ok {
			msg += fmt.Sprintf("%s,", labelKey)
		}
	}
	if msg != "" {
		return false, "Missing labels: " + msg
	}
	return true, ""
}