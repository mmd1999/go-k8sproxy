package main

import (
	"fmt"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
)

type kclient struct {
	client *kubernetes.Clientset
}

func (k kclient) handleDeployments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("%s request received for %s with params %q", r.Method, r.URL.Path, r.URL.Query().Encode())

	ns := r.URL.Query().Get("ns")
	name := r.URL.Query().Get("name")
	verbose := r.URL.Query().Get("verbose")
	resp, err := getDeployments(k.client, ns, name, verbose)
	if err != nil {
		fmt.Fprintf(w, "Error calling %q. %q\n", r.URL.Path, err)
	}
	w.Write(resp)
}

func (k kclient) handleDaemonsets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("%s request received for %s with params %q", r.Method, r.URL.Path, r.URL.Query().Encode())

	ns := r.URL.Query().Get("ns")
	name := r.URL.Query().Get("name")
	verbose := r.URL.Query().Get("verbose")
	resp, err := getDaemonsets(k.client, ns, name, verbose)
	if err != nil {
		fmt.Fprintf(w, "Error calling %q. %q\n", r.URL.Path, err)
	}
	w.Write(resp)
}

func (k kclient) handleCronjobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("%s request received for %s with params %q", r.Method, r.URL.Path, r.URL.Query().Encode())

	ns := r.URL.Query().Get("ns")
	name := r.URL.Query().Get("name")
	verbose := r.URL.Query().Get("verbose")
	resp, err := getCronjobs(k.client, ns, name, verbose)
	if err != nil {
		fmt.Fprintf(w, "Error calling %q. %q\n", r.URL.Path, err)
	}
	w.Write(resp)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
func readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
