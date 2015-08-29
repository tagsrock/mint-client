#! /bin/bash
set -e

# mct = mint client test

cd $GOPATH/src/github.com/eris-ltd/mint-client
echo "********** BUILDING MINT-CLIENT ********"
docker build -t mct_client -f ./DOCKER/linked/DockerfileClient . 

cd ./DOCKER/eris-cli
bash run.sh
