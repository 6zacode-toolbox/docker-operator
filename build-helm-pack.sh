#!/bin/bash
shopt -s expand_aliases
source  ~/.bash_profile

export HERE=$PWD
export CHART_NAME=docker-operator
export CHART_DIR=$PWD/helm/

docker run --rm -it \
    -v $PWD/operator:/go/src \
    --platform linux/amd64 \
    -v $CHART_DIR:/chart \
    -w /go/src \
    6zar/kubebuilder \
    bash -c "kustomize build config/default >  /chart/charts/$CHART_NAME/templates/deploy.yaml"


docker run -it --rm \
    -v $CHART_DIR:/chart \
    --platform linux/amd64 \
    arielev/pybump:1.9.3 \
    bump --file /chart/charts/$CHART_NAME/Chart.yaml --level minor


docker run --rm -it \
    -v $CHART_DIR:/chart \
    --platform linux/amd64 \
    -w /chart \
    6zar/kubebuilder \
    helm package charts/$CHART_NAME/

docker run --rm -it \
    -w /chart \
    --platform linux/amd64 \
    -v $CHART_DIR:/chart \
    6zar/kubebuilder \
    helm repo index --url https://6za.github.io/helm-charts/ --merge index.yaml .