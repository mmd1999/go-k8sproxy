#!/bin/bash

export chart_version=$1
export app_version=$2
echo "chart_version: $chart_version"
echo "app_version: $app_version"

yq -i '.version = env(chart_version)' ./charts/k8sproxy/Chart.yaml
yq -i '.appVersion = strenv(app_version)' ./charts/k8sproxy/Chart.yaml