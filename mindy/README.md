# mindy
Tinydns fed by a blockchain

Commands
--------

Mindy expects tinydns to be installed (in /etc/services/tinydns)

To add all dns entries from the blockchain into tinydns, run

`mindy catchup`

Mindy will run the ListNames rpc call to get all name registry entries,
parse for A-records, and add them to the tinydns data file using the `add-host` or `add-alias` commands of
tinydns as necessary.

To run an active daemon that stays up to date with the blockchain by listening for NameTx events, run

`mindy run`
