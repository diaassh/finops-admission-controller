package main

const (
	CPUCostPerCoreHour   = 2.00
	MemoryCostPerGiBHour = 0.50
	HoursPerMonth        = 730.0
)

type Workload struct {
	CPUCores  float64
	MemoryGiB float64
	Replicas  float64
	Budget    float64
}

func CalculateMonthlyCost(cpuCores, memoryGiB, replicas float64) float64 {
	if cpuCores < 0 {
		cpuCores = 0
	}
	if memoryGiB < 0 {
		memoryGiB = 0
	}
	if replicas < 0 {
		replicas = 0
	}

	cpuCost := cpuCores * CPUCostPerCoreHour * HoursPerMonth
	memoryCost := memoryGiB * MemoryCostPerGiBHour * HoursPerMonth

	return (cpuCost + memoryCost) * replicas
}

func ExceedsBudget(cost, budget float64) bool {
	if cost < 0 {
		cost = 0
	}

	if budget < 0 {
		budget = 0
	}

	return cost > budget
}

func IsValidWorkload(workload Workload) bool {
	if workload.CPUCores <= 0 {
		return false
	}

	if workload.MemoryGiB <= 0 {
		return false
	}

	if workload.Replicas < 1 {
		return false
	}

	if workload.Budget < 0 {
		return false
	}

	return true
}

type EvaluationResult struct {
    EstimatedMonthlyCost float64
    Allowed              bool
    Reason               string
}

func EvaluateWorkload(workload Workload) EvaluationResult {
    if !IsValidWorkload(workload) {
        return EvaluationResult{
            Allowed: false,
            Reason:  "Workload has invalid values.",
        }
    }

    cost := CalculateMonthlyCost(
        workload.CPUCores,
        workload.MemoryGiB,
        workload.Replicas,
    )

    if ExceedsBudget(cost, workload.Budget) {
        return EvaluationResult{
            EstimatedMonthlyCost: cost,
            Allowed:              false,
            Reason:               "Workload exceeds budget.",
        }
    }

    return EvaluationResult{
        EstimatedMonthlyCost: cost,
        Allowed:              true,
        Reason:               "Workload is within budget.",
    }
}

