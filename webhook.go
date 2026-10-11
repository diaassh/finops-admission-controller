
package main

import (
    "encoding/json"
    "net/http"
)

func EvaluateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST requests are allowed",
            http.StatusMethodNotAllowed)
        return
    }

    var workload Workload

    err := json.NewDecoder(r.Body).Decode(&workload)
    if err != nil {
        http.Error(w, "Invalid JSON request body",
            http.StatusBadRequest)
        return
    }

    result := EvaluateWorkload(workload)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
