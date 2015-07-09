#! /bin/sh

cd $GOPATH/src/github.com/eris-ltd/mint-client
docker build -t client -f ./DOCKER/all_in_one/Dockerfile . 
docker run -t --rm client
