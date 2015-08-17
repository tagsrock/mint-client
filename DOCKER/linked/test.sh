#! /bin/bash

# Some simple tests for sending txs through the cli

export MINTX_NODE_ADDR=http://tendermint:46657/
export MINTX_SIGN_ADDR=http://keys:4767
export MINTX_CHAINID=$CHAIN_ID

# check the chain id
STATUS=`mintinfo status`
echo "status $STATUS"
CHAIN_ID2=`echo $STATUS | jq .[1].node_info.chain_id`
CHAIN_ID2=$(echo "$CHAIN_ID2" | tr -d '"') # remove surrounding quotes
echo "chain id $CHAIN_ID2"

if [ "$CHAIN_ID" != "$CHAIN_ID2" ]; then
	echo "Wrong chain id. Got $CHAIN_ID2, expected $CHAIN_ID"
	exit 1
fi

echo "******** RUNNING TEST: NameTx ************"

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

STATUS=`mintinfo status`
echo "status $STATUS"

# verify the name reg entry
DATA=`mintinfo names $REG_NAME data`
echo $DATA

if [ "$REG_DATA" != "$DATA" ]; then
	echo "Wrong data. Got $DATA, expected $REG_DATA"
	exit 1
fi

echo "******** RUNNING TEST: CallTx with wait ************"

# send a calltx that just returns 5
# PUSH1 05 PUSH1 00 MSTORE PUSH1 20 PUSH1 00 RETURN
CODE="600560005260206000F3"
EXPECT="5"
MINTX_OUTPUT=`mintx --debug call --to "" --data $CODE --amt 10 --fee 0 --gas 1000 --sign --broadcast --wait`
EXIT=$?
echo "$MINTX_OUTPUT"
if [ $EXIT -gt 0 ]; then
	echo "Failed to send mint transaction"
	exit 1
fi

RESULT=`echo "$MINTX_OUTPUT" | grep "Return Value:" | awk '{print $3}' | sed 's/^0*//'`
echo "result $RESULT"

if [ "$RESULT" != "$EXPECT" ]; then
	echo "Wrong result. Got $RESULT, expected $EXPECT"
	exit 1
fi

#--------------------------

echo "PASS"
