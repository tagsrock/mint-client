#! /bin/bash

# A simple integrations test wherein we startup a one-man blockchain and register a name on it

# setup a new key
CHAIN_ID=testchain
mintgen random --dir=$TMROOT 1 $CHAIN_ID
ls $TMROOT
cat $TMROOT/genesis.json
ADDR=`mintkey eris $TMROOT/priv_validator.json`
echo "addr $ADDR"
export MINTX_PUBKEY=`eris-keys pub --addr $ADDR`
echo "pub $MINTX_PUBKEY"

# start the daemons
tendermint node &
eris-keys --debug server &

export MINTX_NODE_ADDR=http://localhost:46657/
export MINTX_SIGN_ADDR=http://localhost:4767
export MINTX_CHAINID=$CHAIN_ID

# let tendermint start
sleep 5

# check the chain id
STATUS=`mintinfo status`
echo "status $STATUS"
CHAIN_ID2=`echo $STATUS | jq .chain_id`
echo "chain id $CHAIN_ID2"
CHAIN_ID2=$(echo "$CHAIN_ID2" | tr -d '"') # remove surrounding quotes
echo "chain id $CHAIN_ID2"

if [ "$CHAIN_ID" != "$CHAIN_ID2" ]; then
	echo "Wrong chain id. Got $CHAIN_ID2, expected $CHAIN_ID"
	exit 1
fi

# create a namereg entry
REG_NAME="artifact"
REG_DATA="blue"
mintx --debug name --name $REG_NAME --data $REG_DATA --amt 1000 --fee 0 --sign --broadcast
EXIT=$?
if [ $EXIT -gt 0 ]; then
	echo "Failed to send mint transaction"
	exit 1
fi

sleep 5

# verify the name reg entry
DATA=`mintinfo names $REG_NAME data`
echo $DATA

if [ "$REG_DATA" != "$DATA" ]; then
	echo "Wrong data. Got $DATA, expected $REG_DATA"
	exit 1
fi

echo "PASS"
