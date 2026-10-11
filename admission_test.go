package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
)

func TestKubernetesAdmissionHandler(t *testing.T) {
	body := `{
		"apiVersion": "admission.k8s.io/v1",
		"kind": "AdmissionReview",
		"request": {
			"uid": "test-uid",
			"kind": {
				"group": "apps",
				"version": "v1",
				"kind": "Deployment"
			},
			"object": {
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"metadata": {
					"name": "demo",
					"annotations": {
						"finops.example.com/monthly-budget": "15000"
					}
				},
				"spec": {
					"replicas": 3,
					"template": {
						"spec": {
							"containers": [{
								"name": "app",
								"resources": {
									"requests": {
										"cpu": "2",
										"memory": "4Gi"
									}
								}
							}]
						}
					}
				}
			}
		}
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/validate",
		strings.NewReader(body),
	)
	response := httptest.NewRecorder()

	KubernetesAdmissionHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d", response.Code)
	}

	var review admissionv1.AdmissionReview

	err := json.NewDecoder(response.Body).Decode(&review)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	if review.Response == nil {
		t.Fatal("Expected an AdmissionReview response")
	}

	if !review.Response.Allowed {
		t.Errorf("Expected workload to be allowed, got: %s",
			review.Response.Result.Message)
	}

	if review.Response.UID != "test-uid" {
		t.Error("Response UID does not match request UID")
	}
}

func TestKubernetesAdmissionHandlerRejectsOverBudget(t *testing.T) {
	body := `{
		"apiVersion": "admission.k8s.io/v1",
		"kind": "AdmissionReview",
		"request": {
			"uid": "test-uid",
			"kind": {
				"group": "apps",
				"version": "v1",
				"kind": "Deployment"
			},
			"object": {
				"apiVersion": "apps/v1",
				"kind": "Deployment",
				"metadata": {
					"name": "demo",
					"annotations": {
						"finops.example.com/monthly-budget": "10000"
					}
				},
				"spec": {
					"replicas": 3,
					"template": {
						"spec": {
							"containers": [{
								"name": "app",
								"resources": {
									"requests": {
										"cpu": "2",
										"memory": "4Gi"
									}
								}
							}]
						}
					}
				}
			}
		}
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/validate",
		strings.NewReader(body),
	)
	response := httptest.NewRecorder()

	KubernetesAdmissionHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200, got %d", response.Code)
	}

	var review admissionv1.AdmissionReview

	err := json.NewDecoder(response.Body).Decode(&review)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	if review.Response == nil {
		t.Fatal("Expected an AdmissionReview response")
	}

	if review.Response.Allowed {
		t.Errorf("Expected workload to be rejected because it exceeds the budget")
	}
}
