#! /bin/sh

#------------------------------------------------------------------
# create data containers for tendermint and keys

echo "******** CREATING DATA CONTAINERS ************"
docker run --name mct_tendermint-data eris/data echo "tendermint data container"
docker run --name mct_keys-data eris/data echo "keys data container"

#------------------------------------------------------------------
# setup keys

echo "******** RUNNING KEYS DAEMON ************"
# run eris-keys, generate a key, get pubkey
docker run --name mct_keys --volumes-from mct_keys-data -d -p 4767:4767 mct_keys
ADDR=`docker run -t --rm --volumes-from mct_keys-data mct_keys eris-keys gen`
echo "addr $ADDR"
# ADDR has an extra character. trim it
ADDR=${ADDR%?}
PUBKEY=`docker run -t --rm --volumes-from mct_keys-data mct_keys eris-keys pub --addr $ADDR`
PUBKEY=${PUBKEY%?}
echo "pub $PUBKEY"

# get the TMROOT and set the chain id
TMROOT=`docker run --rm --volumes-from mct_tendermint-data -t mct_tendermint bash -c "mkdir -p \\$TMROOT; echo \\$TMROOT"`
TMROOT=${TMROOT%?}
echo "tmroot $TMROOT"
CHAIN_ID=mintclient_test
echo "chain_id $CHAIN_ID"

# XXX: we need to get the privkey out as priv_validator.json so we can copy into tendermint container
# This step can be eliminated once tendermint can use eris-keys for signing
# NOTE: I tried and failed to do this through an ENV var, so resorted to saving (and removing) priv_validator.json on the host ...
docker run --rm --volumes-from mct_keys-data -t mct_client mintkey mint $ADDR > priv_validator.json
cat priv_validator.json | docker run --rm --volumes-from mct_tendermint-data -i mct_tendermint bash -c "cat > $TMROOT/priv_validator.json"
rm priv_validator.json


#------------------------------------------------------------------
# generate the genesis.json

echo "******** GENERATE GENESIS.JSON ************"

# run mintgen in mint-client with volumes from tendermint-data and using the pubkey
docker run --rm --volumes-from mct_tendermint-data -t mct_client bash -c "mintgen single --pub=$PUBKEY $CHAIN_ID > $TMROOT/genesis.json"
GENESIS=`docker run --rm --volumes-from mct_tendermint-data -t mct_client bash -c "cat $TMROOT/genesis.json"`
echo "genesis $GENESIS"

# copy in the config.toml
cat config.toml | docker run --rm --volumes-from mct_tendermint-data -i mct_tendermint bash -c "cat > $TMROOT/config.toml"

#------------------------------------------------------------------
# start tendermint

echo "******** RUNNING TENDERMINT ************"

# run the tendermint container with volumes from tendermint-data
docker run --name mct_tendermint --volumes-from mct_tendermint-data -d -p 46657:46657 mct_tendermint

# let tendermint start
sleep 3

#------------------------------------------------------------------
# run test

echo "******** RUNNING TEST ************"

# run the test commands in mint-client container linked to eris-keys and tendermint
docker run --name mct_client_test --rm --link mct_keys:keys --link mct_tendermint:tendermint -e "CHAIN_ID=$CHAIN_ID" -e "MINTX_PUBKEY=$PUBKEY" -t mct_client

#------------------------------------------------------------------
# cleanup

echo "-----------"
echo "cleaning up ..."
docker rm -vf mct_keys mct_tendermint mct_keys-data mct_tendermint-data 
