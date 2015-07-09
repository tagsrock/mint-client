#! /bin/sh

#------------------------------------------------------------------
# create data containers for tendermint and keys

docker run --name tendermint-data eris/data echo "tendermint data container"
docker run --name keys-data eris/data echo "keys data container"

#------------------------------------------------------------------
# setup keys

# run eris-keys, generate a key, get pubkey
docker run --name keys --volumes-from keys-data -d -p 4767:4767 keys
ADDR=`docker exec -t keys eris-keys gen`
echo "addr $ADDR"
# ADDR has an extra character. trim it
ADDR=${ADDR%?}
PUBKEY=`docker exec -t keys eris-keys pub --addr $ADDR`
PUBKEY=${PUBKEY%?}
echo "pub $PUBKEY"

# get the TMROOT and set the chain id
TMROOT=`docker run --rm --volumes-from tendermint-data -t tendermint bash -c "mkdir -p \\$TMROOT; echo \\$TMROOT"`
TMROOT=${TMROOT%?}
echo $TMROOT
CHAIN_ID=mintclient_test

# XXX: we need to get the privkey out as priv_validator.json so we can copy into tendermint container
# This step can be eliminated once tendermint can use eris-keys for signing
# NOTE: I tried and failed to do this through an ENV var, so resorted to saving (and removing) priv_validator.json on the host ...
docker run --rm --volumes-from keys-data -t client mintkey mint $ADDR > priv_validator.json
cat priv_validator.json | docker run --rm --volumes-from tendermint-data -i tendermint bash -c "cat > $TMROOT/priv_validator.json"
rm priv_validator.json


#------------------------------------------------------------------
# generate the genesis.json

# run mintgen in mint-client with volumes from tendermint-data and using the pubkey
docker run --rm --volumes-from tendermint-data -t client bash -c "mintgen single --pub=$PUBKEY $CHAIN_ID > $TMROOT/genesis.json"

# copy in the config.toml
cat config.toml | docker run --rm --volumes-from tendermint-data -i tendermint bash -c "cat > $TMROOT/config.toml"

#------------------------------------------------------------------
# start tendermint

# run the tendermint container with volumes from tendermint-data
docker run --name tendermint --volumes-from tendermint-data -d -p 46657:46657 tendermint

# let tendermint start
sleep 3

#------------------------------------------------------------------
# run test

# run the test commands in mint-client container linked to eris-keys and tendermint
docker run --name client_test --rm --link keys:keys --link tendermint:tendermint -e "CHAIN_ID=$CHAIN_ID" -e "MINTX_PUBKEY=$PUBKEY" -t client

#------------------------------------------------------------------
# cleanup

echo "-----------"
echo "cleaning up ..."
docker stop keys tendermint
docker rm keys tendermint keys-data tendermint-data 
