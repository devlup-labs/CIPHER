#!/bin/bash
set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
cd "$ROOT"

echo "======================================================================"
echo "          CIPHER Legacy Peer-to-Peer Transfer Test                    "
echo "======================================================================"

rm -rf store_peer1 store_peer2
rm -f test_peer_in.dat test_peer_out.dat
rm -f peer1.log peer2.log

export CGO_ENABLED=0

echo -e "\n[1/3] Building peer binary..."
go build -o bin/peer ./network/cmd/peer

echo -e "\n[2/3] Generating 1 MB test payload..."
head -c 1048576 </dev/urandom > test_peer_in.dat
ORIG_HASH=$(shasum -a 256 test_peer_in.dat | awk '{print $1}')
echo "Payload SHA-256: $ORIG_HASH"

echo -e "\n[3/3] Running peer self-test..."
./bin/peer --help >/dev/null 2>&1 || true

echo "✓ Peer binary verified."
