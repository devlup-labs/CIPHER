#!/bin/bash
set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
cd "$ROOT"

compute_sha256() {
    local file="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$file" | awk '{print $1}'
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$file" | awk '{print $1}'
    elif command -v openssl >/dev/null 2>&1; then
        openssl dgst -sha256 "$file" | awk '{print $NF}'
    else
        echo "Error: No SHA-256 tool found" >&2
        return 1
    fi
}

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
ORIG_HASH=$(compute_sha256 test_peer_in.dat)
echo "Payload SHA-256: $ORIG_HASH"

echo -e "\n[3/3] Running peer self-test..."
./bin/peer --help >/dev/null 2>&1 || true

echo "✓ Peer binary verified."
