# goproject

A simple application that queries a Kubernetes cluster API and returns JSON output for various workloads.

Supported endpoints
```
/deployments
/daemonsets
/cronjobs
```

Supported query params:
```
ns = namespace
name = workload name
verbose = true to get full manifest output
```

Examples:
```
curl -s "http://localhost:5000/deployments"

curl -s "http://localhost:5000/daemonsets/ns=kube-system&name=kube-proxy"

curl -s "http://localhost:5000/cronjobs/verbose=true"
```