# go-k8sproxy

A simple application that queries a Kubernetes cluster API and returns JSON output for various workloads.

Supports running as an application in-cluster and out-of-cluster

#### Environment vars:
| Environment Var | Default | Description |
| ----------------|---------|-------------|
| IN_CLUSTER | `true` | set to false to run application out-of-cluster |
| KUBE_CONFIG_PATH | `$HOME/.kube/config` | used for out-of-cluster mode only `/<path>/<to>/<your>/<kubeconfig>/<file>` |

#### Supported endpoints:
```
/deployments
/daemonsets
/cronjobs
```

#### Supported query params:
```
ns = namespace
name = workload name
verbose = true to get full manifest output
```

#### Examples:
```
curl -s "http://localhost:5000/deployments"

curl -s "http://localhost:5000/daemonsets/ns=kube-system&name=kube-proxy"

curl -s "http://localhost:5000/cronjobs/verbose=true"
```

#### Example output:
```
curl -s "http://localhost:5000/deployments"
[
  {
    "name": "coredns",
    "namespace": "kube-system",
    "Container": [
      {
        "imagename": "registry.k8s.io/coredns/coredns",
        "shortname": "coredns",
        "version": "v1.11.3"
      }
    ]
  },
  {
    "name": "local-path-provisioner",
    "namespace": "local-path-storage",
    "Container": [
      {
        "imagename": "docker.io/kindest/local-path-provisioner",
        "shortname": "local-path-provisioner",
        "version": "v20250214-acbabc1a"
      }
    ]
  }
]
```