#! /bin/bash

# A dead simple test to make sure we can sign a transaction through the mintx cli

ADDR=`eris-keys gen`
echo "addr $ADDR"
export MINTX_PUBKEY=`eris-keys pub --addr $ADDR`
echo "pub $MINTX_PUBKEY"
eris-keys server &
sleep 1


mintx name --verbose --name name --data data --amt 100 --fee 0 --nonce 1 --sign
EXIT=$?

if [ $EXIT -gt 0 ]; then
	echo "FAIL!"
else
	echo "PASS"
fi


