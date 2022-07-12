#!/bin/sh

build_push(){
  docker buildx build  --platform ${ARCHS} -t ${REGISTRY}/${NAME}:latest  --output=type=registry,registry.insecure=true --push .
}

helm_build_push(){
  FN=${NAME}-${VER}.tgz
  rm ${FN}
  helm package ./install --version ${VER}
  curl --data-binary "@${FN}" http://helm.local/api/charts
}

REGISTRY=registry.local
NAME=hs-router
ARCHS="linux/amd64,linux/arm64,linux/386"
VER=0.1.10


#helm_build_push
build_push





