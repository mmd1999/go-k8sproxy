package main

import (
	"context"
	"encoding/json"
	"strings"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// build a custom struct to marshal deployment data into
type Cronjob struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container []Container
}

func getCronjobs(k *kubernetes.Clientset, ns string, name string, full string) ([]byte, error) {
	cjl := &v1.CronJobList{}
	var err error
	if name != "" {
		c, err := k.BatchV1().CronJobs(ns).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		cjl.Items = append(cjl.Items, *c)
	} else {
		cjl, err = k.BatchV1().CronJobs(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	}
	var output []byte
	if full != "true" {
		var cronjobs []Cronjob
		for _, cj := range cjl.Items {
			var container []Container
			for _, v := range cj.Spec.JobTemplate.Spec.Template.Spec.Containers {
				ss := strings.Split(v.Image, ":")
				splitName := strings.Split(ss[0], "/")
				container = append(container, Container{
					ImageNameFull:  ss[0],
					ImageNameShort: splitName[len(splitName)-1],
					Version:        ss[1],
				})
			}
			cronjobs = append(cronjobs, Cronjob{
				Name:      cj.Name,
				Namespace: cj.Namespace,
				Container: container,
			})
		}
		output, err = json.MarshalIndent(cronjobs, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		output, err = json.MarshalIndent(cjl, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}
