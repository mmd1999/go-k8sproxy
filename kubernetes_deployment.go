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
type Deployment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container []Container
}

type Container struct {
	ImageNameFull  string `json:"imagename"`
	ImageNameShort string `json:"shortname"`
	Version        string `json:"version"`
}

func getDeployments(k *kubernetes.Clientset, ns string, name string, full string) ([]byte, error) {
	dl := &v1.DeploymentList{}
	var err error
	if name != "" {
		d, err := k.AppsV1().Deployments(ns).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		dl.Items = append(dl.Items, *d)
	} else {
		dl, err = k.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	}
	var output []byte
	if full != "true" {
		var deployments []Deployment
		for _, deployment := range dl.Items {
			var container []Container
			for _, v := range deployment.Spec.Template.Spec.Containers {
				ss := strings.Split(v.Image, ":")
				splitName := strings.Split(ss[0], "/")
				container = append(container, Container{
					ImageNameFull:  ss[0],
					ImageNameShort: splitName[len(splitName)-1],
					Version:        ss[1],
				})
			}
			deployments = append(deployments, Deployment{
				Name:      deployment.Name,
				Namespace: deployment.Namespace,
				Container: container,
			})
		}
		output, err = json.MarshalIndent(deployments, "", "  ")
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

// working, marshals just the fields we want to json output
// type DeploymentList struct {
// 	Kind       string       `json:"kind"`
// 	APIVersion string       `json:"apiVersion"`
// 	Items      []Deployment `json:"items"`
// }

// type Deployment struct {
// 	Metadata struct {
// 		Name      string `json:"name"`
// 		Namespace string `json:"namespace"`
// 	} `json:"metadata"`
// 	Spec struct {
// 		Template struct {
// 			Spec struct {
// 				Containers []struct {
// 					Image string `json:"image"`
// 				} `json:"containers"`
// 			} `json:"spec"`
// 		} `json:"template"`
// 	} `json:"spec"`
// }

// func getDeployments(k *kubernetes.Clientset) ([]byte, error) {
// 	objs, err := k.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		fmt.Printf("Unable to read Deployment objects from cluster: %s", err)
// 		return nil, err
// 	}
// 	deployList, err := json.Marshal(objs)
// 	if err != nil {
// 		fmt.Printf("Unable to mmarshal json for Deployments: %s", err)
// 		return nil, err
// 	}
// 	var dl DeploymentList
// 	err = json.Unmarshal(deployList, &dl)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	json, err := json.Marshal(dl.Items)
// 	return json, nil
// }
