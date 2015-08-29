# Docker and Tests

Here we have some integrations tests that run in/across docker containers.

The basic test is to create a key, start a blockchain with one validator (that key), 
send a name reg transaction to the chain with some data, and verify that data exists on the chain.
We also deploy a contract that just returns something.

This feat makes partial use of mintgen, mintkey, mintinfo, and mintx.

`all-in-one` contains the test implemented in a single docker container. It is deprecated and probably doesn't work.

`linked` contains the test implemented across three docker containers, one for each of tendermint, eris-keys, and the mint-client. 
There is lots of docker plumbing

`eris-cli` contains the test implemented using eris-cli to create the chain and manage containers.

Each folder has a build script (that will build containers, configure/start them and run the test), a run script (configure/start containers and run the test), and a test script (the actual transaction tests running in a client_test container). 

NOTE: `test.sh` is (must be) replicated in each folder because Docker, in its infinite wisdom, doesn't allow symlinks.

