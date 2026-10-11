package main

import "testing"

func TestCalculateMonthlyCost(t *testing.T) {
	cost := CalculateMonthlyCost(2, 4, 3)

	expected := 13140.0

	if cost != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, cost)
	}
}

func TestExceedsBudget(t *testing.T) {
	if !ExceedsBudget(13140, 10000) {
		t.Error("Expected workload to exceed budget")
	}

	if ExceedsBudget(5000, 15000) {
		t.Error("Expected workload to stay within budget")
	}
}

func TestIsValidWorkload(t *testing.T) {
	valid := Workload{
		CPUCores:  2.0,
		MemoryGiB: 4.0,
		Replicas:  3.0,
		Budget:    15000.0,
	}

	if !IsValidWorkload(valid) {
		t.Error("Expected workload to be valid")
	}

	invalid := valid
	invalid.CPUCores = -1

	if IsValidWorkload(invalid) {
		t.Error("Expected workload with negative CPU to be invalid")
	}
}
func TestEvaluateWorkload(t *testing.T) {
    // Case 1: Valid workload within budget
    workload := Workload{
        CPUCores:  2.0,
        MemoryGiB: 4.0,
        Replicas:  3.0,
        Budget:    15000.0,
    }

    result := EvaluateWorkload(workload)

    if !result.Allowed {
        t.Error("Expected workload to be allowed")
    }

    if result.EstimatedMonthlyCost != 13140.0 {
        t.Errorf("Expected cost 13140, got %.2f",
            result.EstimatedMonthlyCost)
    }

    // Case 2: Workload exceeds budget
    workload.Budget = 10000.0
    result = EvaluateWorkload(workload)

    if result.Allowed {
        t.Error("Expected over-budget workload to be rejected")
    }

    // Case 3: Invalid workload
    workload.CPUCores = -1.0
    result = EvaluateWorkload(workload)

    if result.Allowed {
        t.Error("Expected invalid workload to be rejected")
    }
}
