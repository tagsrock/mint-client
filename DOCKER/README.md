# Docker and Tests

Here we have some integrations tests that run in/across docker containers.

The basic test is to create a key, start a blockchain with one validator (that key), 
send a name reg transaction to the chain with some data, and verify that data exists on the chain.

This feat makes partial use of mintgen, mintkey, mintinfo, and mintx.

`all-in-one` contains the test implemented in a single docker container

`linked` contains the test implemented across three docker containers, one for each of tendermint, eris-keys, and the mint-client
