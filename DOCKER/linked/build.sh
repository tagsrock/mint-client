#! /bin/sh
set -e

cd $GOPATH/src/github.com/eris-ltd/mint-client
echo "********** BUILDING TENDERMINT ********"
docker build -t tendermint -f ./DOCKER/linked/DockerfileTendermint . 
echo "********** BUILDING ERIS-KEYS ********"
docker build -t keys -f ./DOCKER/linked/DockerfileKeys . 
echo "********** BUILDING MINT-CLIENT ********"
docker build -t client -f ./DOCKER/linked/DockerfileClient . 

cd ./DOCKER/linked
./run.sh
