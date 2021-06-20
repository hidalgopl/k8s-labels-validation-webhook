package main

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestHasRecommendedLabels(t *testing.T) {
	cases := []struct {
		name    string
		pod     v1.Pod
		allowed bool
		msg     string
	}{
		{
			name: "missing labels",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "pod1",
					Labels: map[string]string{},
				},
			},
			allowed: false,
			msg:     "Missing labels: app.kubernetes.io/name,app.kubernetes.io/instance,app.kubernetes.io/version,app.kubernetes.io/component,app.kubernetes.io/part-of,app.kubernetes.io/managed-by,app.kubernetes.io/created-by,",
		},
		{
			name: "has all labels set",
			pod: v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod1",
					Labels: map[string]string{
						"app.kubernetes.io/name":       "dummy",
						"app.kubernetes.io/instance":   "dummy",
						"app.kubernetes.io/version":    "dummy",
						"app.kubernetes.io/component":  "dummy",
						"app.kubernetes.io/part-of":    "dummy",
						"app.kubernetes.io/managed-by": "dummy",
						"app.kubernetes.io/created-by": "dummy",
					},
				},
			},
			allowed: true,
			msg:     "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			isAllowed, msg := hasRecommendedLabels(&tc.pod)
			assert.Equal(t, tc.allowed, isAllowed)
			assert.Equal(t, tc.msg, msg)
		})
	}
}
