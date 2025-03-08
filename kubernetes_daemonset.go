package main

import (
	"context"
	"encoding/json"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// build a custom struct to marshal deployment data into
type Daemonset struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container []Container
}

func getDaemonsets(k *kubernetes.Clientset, ns string, name string, full string) ([]byte, error) {
	dl := &v1.DaemonSetList{}
	var err error
	if name != "" {
		d, err := k.AppsV1().DaemonSets(ns).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		dl.Items = append(dl.Items, *d)
	} else {
		dl, err = k.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	}
	var output []byte
	if full != "true" {
		var daemonsets []Daemonset
		for _, daemonset := range dl.Items {
			var container []Container
			for _, v := range daemonset.Spec.Template.Spec.Containers {
				ss := strings.Split(v.Image, ":")
				splitName := strings.Split(ss[0], "/")
				container = append(container, Container{
					ImageNameFull:  ss[0],
					ImageNameShort: splitName[len(splitName)-1],
					Version:        ss[1],
				})
			}
			daemonsets = append(daemonsets, Daemonset{
				Name:      daemonset.Name,
				Namespace: daemonset.Namespace,
				Container: container,
			})
		}
		output, err = json.MarshalIndent(daemonsets, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		output, err = json.MarshalIndent(dl, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}
