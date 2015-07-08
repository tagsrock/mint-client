#! /bin/bash

# A simple test to make sure we can send a transaction through the mintx cli

export MINTX_NODE_ADDR=http://tendermint:46657/
export MINTX_SIGN_ADDR=http://keys:4767
export MINTX_CHAINID=$CHAIN_ID

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
