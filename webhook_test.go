
package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestEvaluateHandler(t *testing.T) {
    body := `{
        "CPUCores": 2,
        "MemoryGiB": 4,
        "Replicas": 3,
        "Budget": 15000
    }`

    request := httptest.NewRequest(
        http.MethodPost,
        "/evaluate",
        strings.NewReader(body),
    )

    response := httptest.NewRecorder()

    EvaluateHandler(response, request)

    if response.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", response.Code)
    }

    var result EvaluationResult

    err := json.NewDecoder(response.Body).Decode(&result)
    if err != nil {
        t.Fatalf("Could not read response: %v", err)
    }

    if !result.Allowed {
        t.Error("Expected workload to be allowed")
    }

    if result.EstimatedMonthlyCost != 13140 {
        t.Errorf("Expected cost 13140, got %.2f",
            result.EstimatedMonthlyCost)
    }
}
