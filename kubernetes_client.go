package main

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Setups up kubernetes client authentication
func (c config) AuthK8s() (*kubernetes.Clientset, error) {
	var conf *rest.Config
	var err error
	// creates an in-cluster config
	if c.InCluster {
		log.Println("Setting up in-cluster authentication")
		conf, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// creates an out-of-cluster config
		var kubeconfig string
		log.Println("Setting up out-of-cluster authentication")
		if c.KubeConfigPath != "" {
			kubeconfig = c.KubeConfigPath
		} else if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			kubeconfig = ""
		}
		log.Printf("Kubernetes config path set to: %s", kubeconfig)
		// use the current context in kubeconfig
		conf, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		log.Printf("Kubernetes host: %s", conf.Host)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
