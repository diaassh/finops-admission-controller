package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const MonthlyBudgetAnnotation = "finops.example.com/monthly-budget"

func KubernetesAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed",
			http.StatusMethodNotAllowed)
		return
	}

	var review admissionv1.AdmissionReview

	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid AdmissionReview request",
			http.StatusBadRequest)
		return
	}

	if review.Request == nil {
		http.Error(w, "AdmissionReview request is missing",
			http.StatusBadRequest)
		return
	}

	request := review.Request

	response := &admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: false,
	}

	deny := func(reason string) {
		response.Allowed = false
		response.Result = &metav1.Status{
			Message: reason,
		}
		writeAdmissionReview(w, &review, response)
	}

	if request.Kind.Group != "apps" ||
		request.Kind.Version != "v1" ||
		request.Kind.Kind != "Deployment" {
		deny("Only apps/v1 Deployment resources are supported.")
		return
	}

	var deployment appsv1.Deployment

	if err := json.Unmarshal(request.Object.Raw, &deployment); err != nil {
		deny("Could not read the Deployment object.")
		return
	}

	budgetText, found := deployment.Annotations[MonthlyBudgetAnnotation]
	if !found {
		deny("Missing monthly budget annotation: " +
			MonthlyBudgetAnnotation)
		return
	}

	budget, err := strconv.ParseFloat(strings.TrimSpace(budgetText), 64)
	if err != nil || math.IsNaN(budget) ||
		math.IsInf(budget, 0) || budget < 0 {
		deny("Monthly budget must be a valid, non-negative number.")
		return
	}

	containers := deployment.Spec.Template.Spec.Containers

	if len(containers) == 0 {
		deny("Deployment must contain at least one application container.")
		return
	}

	var totalCPUCores float64
	var totalMemoryGiB float64

	for _, container := range containers {
		cpu, hasCPU := container.Resources.Requests[corev1.ResourceCPU]
		if !hasCPU || cpu.MilliValue() <= 0 {
			deny(fmt.Sprintf(
				"Container %q must have a positive CPU request.",
				container.Name,
			))
			return
		}

		memory, hasMemory := container.Resources.Requests[corev1.ResourceMemory]
		if !hasMemory || memory.Value() <= 0 {
			deny(fmt.Sprintf(
				"Container %q must have a positive memory request.",
				container.Name,
			))
			return
		}

		totalCPUCores += float64(cpu.MilliValue()) / 1000.0
		totalMemoryGiB += float64(memory.Value()) /
			float64(1024*1024*1024)
	}

	replicas := 1.0
	if deployment.Spec.Replicas != nil {
		replicas = float64(*deployment.Spec.Replicas)
	}

	workload := Workload{
		CPUCores:  totalCPUCores,
		MemoryGiB: totalMemoryGiB,
		Replicas:  replicas,
		Budget:    budget,
	}

	result := EvaluateWorkload(workload)

	message := result.Reason
	if result.EstimatedMonthlyCost > 0 {
		message = fmt.Sprintf(
			"%s Estimated monthly cost: ₹%.2f; budget: ₹%.2f.",
			result.Reason,
			result.EstimatedMonthlyCost,
			budget,
		)
	}

	response.Allowed = result.Allowed
	response.Result = &metav1.Status{
		Message: message,
	}

	writeAdmissionReview(w, &review, response)
}

func writeAdmissionReview(
	w http.ResponseWriter,
	review *admissionv1.AdmissionReview,
	response *admissionv1.AdmissionResponse,
) {
	review.Response = response
	review.Request = nil

	review.TypeMeta = metav1.TypeMeta{
		APIVersion: "admission.k8s.io/v1",
		Kind:       "AdmissionReview",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(review)
}
