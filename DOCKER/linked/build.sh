#! /bin/sh
set -e

cd $GOPATH/src/github.com/eris-ltd/mint-client
docker build -t tendermint -f ./DOCKER/linked/DockerfileTendermint . 
docker build -t keys -f ./DOCKER/linked/DockerfileKeys . 
docker build -t client -f ./DOCKER/linked/DockerfileClient . 
