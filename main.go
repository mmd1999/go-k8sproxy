package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
)

const port = ":5000"

type config struct {
	InCluster      bool   `env:"IN_CLUSTER" envDefault:"true"`
	KubeConfigPath string `env:"KUBE_CONFIG_PATH"`
}

func main() {

	// Setup k8s auth and configure client
	var cfg config
	var k kclient
	if err := env.Parse(&cfg); err != nil {
		fmt.Println(err)
	}
	var err error
	k.client, err = cfg.AuthK8s()
	if err != nil {
		log.Fatal(err)
	}

	// Start http Server
	fmt.Printf("Server listening on port %s\n", port)
	mux := http.NewServeMux()

	// HTTP handlers
	mux.HandleFunc("/deployments", k.handleResource)
	mux.HandleFunc("/daemonsets", k.handleResource)
	mux.HandleFunc("/cronjobs", k.handleResource)

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/readiness", readiness)
	log.Fatal(http.ListenAndServe(port, mux))
}
