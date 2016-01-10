#! /bin/bash

ifExit(){
	if [ $? -ne 0 ]; then
		echo "ifExit: $1"
		exit 1
	fi
}


#default
NUM_NODES=1

# XXX: cleanup at beginning regardless
if [ "$CLEANUP" == "true" ]; then
for ((i=1; i<=$NUM_NODES; i++ ))
do	
docker rm -vf mct_keys"_"$i mct_tendermint"_"$i mct_keys-data"_"$i mct_tendermint-data"_"$i
done
exit
fi

#------------------------------------------------------------------
# create data containers for tendermint and keys

echo "******** CREATING DATA CONTAINERS ************"

for ((i=1; i<=$NUM_NODES; i++ ))
do
docker run --name mct_tendermint-data"_"$i eris/data echo "tendermint data container_$i"
docker run --name mct_keys-data"_"$i eris/data echo "keys data container_$i"
done

#------------------------------------------------------------------
# setup keys

echo "******** RUNNING KEYS DAEMON ************"
# run eris-keys, generate a key, get pubkey

for ((i=1; i<=$NUM_NODES; i++)) 
do
docker run --name mct_keys"_"$i  --volumes-from mct_keys-data"_"$i -d mct_keys
ifExit "failed to start keys"
ADDR=`docker exec mct_keys"_"$i  eris-keys gen --no-pass`
# ADDR has an extra character. trim it
ADDRS[$i]=${ADDR}
echo "addr_$i ${ADDRS[$i]}"
PUBKEY=`docker exec mct_keys"_"$i eris-keys pub --addr=${ADDRS[$i]}`
PUBKEYS[$i]=${PUBKEY}
echo "pubkey_$i ${PUBKEYS[$i]}"
done


# get the TMROOT and set the chain id

for ((i=1; i<=$NUM_NODES; i++ ))
do
  TMROOT=`docker run --rm --volumes-from mct_tendermint-data"_"$i -t mct_tendermint bash -c "mkdir -p \\$TMROOT; echo \\$TMROOT"`
done
TMROOT=${TMROOT%?}
echo "tmroot $TMROOT"
CHAIN_ID=mintclient_test
echo "chain_id $CHAIN_ID"

# XXX: we need to get the privkey out as priv_validator.json so we can copy into tendermint container
# This step can be eliminated once tendermint can use eris-keys for signing
# NOTE: I tried and failed to do this through an ENV var, so resorted to saving (and removing) priv_validator.json on the host ...
for ((i=1; i<=$NUM_NODES; i++ ))
do
  docker run --rm --volumes-from mct_keys-data"_"$i -t mct_client mintkey mint ${ADDRS[$i]} > priv_validator.json
  cat priv_validator.json | docker run --rm --volumes-from mct_tendermint-data"_"$i -i mct_tendermint bash -c "cat > $TMROOT/priv_validator.json"
  rm priv_validator.json
done

#------------------------------------------------------------------
# generate the genesis.json

echo "******** GENERATE GENESIS.JSON ************"

# run mintgen in mint-client with volumes from tendermint-data and using the pubkey
for ((i=1; i<=$NUM_NODES; i++ ))
do	

pks="${PUBKEYS[@]}"
docker run --rm --volumes-from mct_tendermint-data"_"$i -t mct_client bash -c "mintgen known --pub=\"$pks\" $CHAIN_ID > $TMROOT/genesis.json"
GENESIS=`docker run --rm --volumes-from mct_tendermint-data"_"$i -t mct_client bash -c "cat $TMROOT/genesis.json"`
done
echo "genesis $GENESIS"

# make the config.toml

docker run --rm --entrypoint mintconfig mct_client --moniker="test_mon" | docker run --rm --volumes-from mct_tendermint-data_1 -i mct_tendermint bash -c "cat > $TMROOT/config.toml"
for ((i=2; i<=$NUM_NODES; i++ ))
do	
docker run --rm --entrypoint mintconfig mct_client --moniker="test_nom" --seeds="mct_tendermint_1:46656"  | docker run --rm --volumes-from mct_tendermint-data"_"$i -i mct_tendermint bash -c "cat > $TMROOT/config.toml"
done
#------------------------------------------------------------------
# start tendermint

echo "******** RUNNING TENDERMINT ************"
echo "starting seed node..."
#set to the first cont
docker run --name mct_tendermint_1 --volumes-from mct_tendermint-data_1 -d mct_tendermint

# run the tendermint container with volumes from tendermint-data
# start i at two, linking each cont to the first, seed is set in config file
for ((i=2; i<=$NUM_NODES; i++)) 
do	
docker run --name mct_tendermint"_"$i --volumes-from mct_tendermint-data"_"$i --link mct_tendermint_1:tendermint -d mct_tendermint 
done
# let tendermint start
sleep 3

#------------------------------------------------------------------
# run test

echo "******** INITIALIZING TESTS ************"
for ((i=1; i<=$NUM_NODES; i++ ))
do	
# run the test commands in mint-client container linked to eris-keys and tendermint
docker run --name mct_client_test"_"$i -d --link mct_keys"_"$i:keys --link mct_tendermint"_"$i:tendermint -e "CHAIN_ID=$CHAIN_ID" -e "PUBKEY=${PUBKEYS[$i]}" -e "I=$i" -e "NUM_NODES=$NUM_NODES" -t mct_client
done

#stdout logs
for ((i=1; i<=$NUM_NODES; i++ ))
do	
docker logs --follow mct_client_test"_"$i 
done
#------------------------------------------------------------------
# cleanup

if [ "$DONT_CLEANUP" != "true" ]; then
	if [ "$CIRCLECI" != "true" ]; then
    	    echo "-----------"
  	    echo "cleaning up ..."
  	    for ((i=1; i<=$NUM_NODES; i++ ))
  	    do	
  	    docker rm -vf mct_keys"_"$i mct_tendermint"_"$i mct_keys-data"_"$i mct_tendermint-data"_"$i mct_client_test"_"$i
  	    done
        fi
fi
