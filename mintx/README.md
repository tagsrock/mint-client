# mintx
Low-level client for talking to tendermint chains

Example
-------
```
mintx name --chainID tendermint_testnet_5e --pubkey F6C79CF0CB9D66B677988BCB9B8EADD9A091CD465A60542A8AB85476256DBA92 --amt 1000 --fee 20 --nonce 2 --name casey --data "psh, we're lawyers, don't tell us how to incorporate" --sign --sign-addr http://localhost:4767 --broadcast --node-addr http://pinkpenguin.chaintest.net:46657 
```

You can simplify by setting some env vars:

```
export MINTX_SIGN_ADDR=http://localhost:4767
export MINTX_NODE_ADDR=http://pinkpenguin.chaintest.net:46657/
export MINTX_PUBKEY=F6C79CF0CB9D66B677988BCB9B8EADD9A091CD465A60542A8AB85476256DBA92
export MINTX_CHAINID=tendermint_testnet_5e
mintx name -amt 1000 --fee 20 --name casey --data "psh, we're lawyers, don't tell us how to incorporate" --sign --broadcast
```

If you don't provide a nonce, and the NODE_ADDR is set, it will fetch the correct nonce for you.
