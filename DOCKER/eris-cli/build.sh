#! /bin/bash
set -e

# mct = mint client test

REPO=github.com/eris-ltd/mint-client
if [[ "$ERIS_VERSION" == "" ]]; then
	ERIS_VERSION=0.11
fi

cd $GOPATH/src/github.com/eris-ltd/mint-client
echo "********** BUILDING MINT-CLIENT ********"
docker build -t mct_client -f ./DOCKER/linked/DockerfileClient . 

docker run --name eris-data eris/data echo "Data-container for testing with eris-cli"

docker cp $GOPATH/src/$REPO/ eris-data:/home/eris/.eris/mint-client/

# eris/eris container in which we run the tests
docker run -t --rm --volumes-from eris-data -v /var/run/docker.sock:/var/run/docker.sock --entrypoint bash quay.io/eris/eris:$ERIS_VERSION /home/eris/.eris/mint-client/DOCKER/eris-cli/run.sh
