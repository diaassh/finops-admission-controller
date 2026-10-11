package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/validate", KubernetesAdmissionHandler)
	
	http.HandleFunc("/", DashboardHandler)

	http.HandleFunc("/evaluate", EvaluateHandler)

	log.Println("FinOps admission server starting with HTTPS on :8443")

	err := http.ListenAndServeTLS(
		":8443",
		"certs/tls.crt",
		"certs/tls.key",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}
