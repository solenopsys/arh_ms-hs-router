#!/bin/sh

build_push(){
  nerdctl build --platform=${ARCHS} --output type=image,name=${REGISTRY}/${NAME}:latest,push=true .
}

helm_build_push(){
  FN=${NAME}-${VER}.tgz
  helm package ./install --version ${VER}
  curl --data-binary "@${FN}" http://helm.klogsolenopsys.org/api/charts
  rm ${FN}
}

REGISTRY=registry.klogsolenopsys.org
NAME=alexstorm-hsm-router
ARCHS="linux/amd64,linux/arm64"
VER=0.1.24


helm_build_push
build_push





