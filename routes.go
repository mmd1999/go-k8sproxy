package main

import (
	"fmt"
	"log"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"k8s.io/client-go/kubernetes"
)

type kclient struct {
	client *kubernetes.Clientset
}

func (k kclient) handleResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Request received for %s with params %q", r.URL.Path, r.URL.Query().Encode())

	ns := r.URL.Query().Get("ns")
	name := r.URL.Query().Get("name")
	verbose := r.URL.Query().Get("verbose")

	err := validation(ns, name, verbose)
	if err != nil {
		log.Printf("Error processing request for %s: %v", r.URL.Path, err)
		http.Error(w, fmt.Sprintf("Invalid query parameters: %v", err), http.StatusBadRequest)
		return
	}

	var resp []byte
	switch r.URL.Path {
	case "/deployments":
		resp, err = getDeployments(k.client, ns, name, verbose)
	case "/daemonsets":
		resp, err = getDaemonsets(k.client, ns, name, verbose)
	case "/cronjobs":
		resp, err = getCronjobs(k.client, ns, name, verbose)
	}
	if err != nil {
		log.Printf("Error processing request for %s: %v", r.URL.Path, err)
		http.Error(w, fmt.Sprintf("Error calling %q: %v", r.URL.Path, err), http.StatusInternalServerError)
		return
	}

	log.Printf("Processed request for %s with params %q", r.URL.Path, r.URL.Query().Encode())
	w.Header().Set("Content-Type", "application/json")

	if string(resp) == "null" {
		http.Error(w, "[]", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
func readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ready"}`))
}

func validation(ns string, name string, verbose string) error {
	v := validator.New()
	// Validate query params if provided
	if name != "" {
		err := v.Var(name, "hostname_rfc1123")
		if err != nil {
			return fmt.Errorf("Invalid 'name' provided, does not follow RFC 1123 format")
		}
	}
	if ns != "" {
		err := v.Var(ns, "hostname_rfc1123")
		if err != nil {
			return fmt.Errorf("Invalid 'namespace' provided, does not follow RFC 1123 format")
		}
	}
	if verbose != "" {
		err := v.Var(verbose, "boolean")
		if err != nil {
			return fmt.Errorf("Invalid value for 'verbose' provided, one of either 'true' or 'false' is expected")
		}
	}
	return nil
}
