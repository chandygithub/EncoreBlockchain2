#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
LANGUAGE=${1:-"golang"}
# CC_SRC_PATH=github.com/

#if [ "$LANGUAGE" = "node" -o "$LANGUAGE" = "NODE" ]; then
#	CC_SRC_PATH=/opt/gopath/src/github.com/fabcar/node
#fi

# clean the keystore
rm -rf ./hfc-key-store

# launch network; create channel and join peer to channel
cd ./first-network
./byfn.sh -m generate -c myc

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars
docker-compose -f docker-compose-cli.yaml up -d

docker exec -it cli bash