# mintinfo
Low-level client for making rpc calls to tendermint chains

# Examples

```
$ mintinfo status
{
        "moniker": "pinkpenguin",
        "chain_id": "tendermint_testnet_5e",
        "version": "0.3.0",
        "genesis_hash": "52063A04D28183292F654257FFE19687D9D9C921CE23E70124AC25D2E45066D7",
        "pub_key": [
                1,
                "F6C79CF0CB9D66B677988BCB9B8EADD9A091CD465A60542A8AB85476256DBA92"
        ],
        "latest_block_hash": "C9EA5B6C085746A829BCDB7B9002A6089D27EB425603A2B140AAC734C722033F",
        "latest_block_height": 14649,
        "latest_block_time": 1433376819383640576
}

$ mintinfo status genesis_hash
"52063A04D28183292F654257FFE19687D9D9C921CE23E70124AC25D2E45066D7"

$ mintinfo status chain_id
"tendermint_testnet_5e"
```

# Env Vars

```
export MINTX_NODE_ADDR=http://pinkpenguin.chaintest.net:46657/
```

or use the `--node-addr` flag 

```
mintinfo --node-addr http://localhost:46657 status
```
