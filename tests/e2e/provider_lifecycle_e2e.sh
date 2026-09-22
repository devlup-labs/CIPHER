#!/bin/bash
set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
cd "$ROOT"

echo "======================================================================"
echo "    CIPHER Lifecycle Test: Provider Independence & Persistence       "
echo "======================================================================"

rm -rf store_provider store_client1 store_client2
rm -f test_orig.dat test_recov1.dat test_recov2.dat
rm -f bootstrap.log publisher.log provider.log client1.log client2.log

export CGO_ENABLED=0

cleanup() {
    kill $BOOT_PID $PUB_PID $PROV_PID 2>/dev/null || true
}
trap cleanup EXIT

echo -e "\n[Step 1/6] Building binaries..."
go build -o bin/bootstrap ./network/cmd/bootstrap
go build -o bin/publisher ./nodes/publisher
go build -o bin/provider ./nodes/provider
go build -o bin/consumer ./nodes/consumer
cp bin/consumer bin/client

echo -e "\n[Step 2/6] Generating 2 MB test payload..."
head -c 2097152 </dev/urandom > test_orig.dat
ORIG_HASH=$(shasum -a 256 test_orig.dat | awk '{print $1}')
echo "Original Payload SHA-256: $ORIG_HASH"

echo -e "\n[Step 3/6] Starting DHT Bootstrap Node..."
./bin/bootstrap -p 48001 -ws-port 0 -identity ./store_provider/boot.key > bootstrap.log 2>&1 &
BOOT_PID=$!
sleep 1

BOOT_ADDR=$(grep "127.0.0.1/tcp/48001/p2p/" bootstrap.log | head -n 1 | awk '{print $NF}')
echo "Bootstrap Multiaddr: $BOOT_ADDR"

echo -e "\n[Step 4/6] Publisher ingests content into Provider store and EXITS..."
./bin/publisher -seed=false -identity ./store_provider/pub.key -store ./store_provider -file test_orig.dat > publisher.log 2>&1

CONTENT_ID=$(grep "^ContentID" publisher.log | awk '{print $NF}')
KEY=$(grep "^Decryption Key" publisher.log | awk '{print $NF}')

echo "Content Published:"
echo "  - ContentID:      $CONTENT_ID"
echo "  - Decryption Key: $KEY"

echo "✓ Verified: Publisher process completed and exited."

echo -e "\n[Step 5/6] Starting Standalone Provider..."
./bin/provider -p 48010 -ws-port 48011 -identity ./store_provider/prov.key -store ./store_provider -bootstrap "$BOOT_ADDR" > provider.log 2>&1 &
PROV_PID=$!
sleep 2

PROV_ID=$(grep "Provider Peer ID:" provider.log | awk '{print $NF}')
echo "Provider running with Peer ID: $PROV_ID"

echo -e "\n[Step 6/6] Consumer 1 retrieves content via DHT without Publisher online..."
./bin/consumer -fetch "$CONTENT_ID" -key "$KEY" -out test_recov1.dat -bootstrap "$BOOT_ADDR" -store ./store_client1 > client1.log 2>&1

RECOV1_HASH=$(shasum -a 256 test_recov1.dat | awk '{print $1}')
echo "Consumer 1 Downloaded SHA-256: $RECOV1_HASH"

if [ "$ORIG_HASH" != "$RECOV1_HASH" ]; then
    echo "ERROR: Hash mismatch on Consumer 1 recovery!"
    cat client1.log
    exit 1
fi
echo "✓ SUCCESS: Content retrieved purely from Provider via DHT!"

echo -e "\n[Persistence Test] Restarting Provider node..."
kill $PROV_PID
sleep 1

./bin/provider -p 48010 -ws-port 48011 -identity ./store_provider/prov.key -store ./store_provider -bootstrap "$BOOT_ADDR" >> provider.log 2>&1 &
PROV_PID=$!
sleep 2
echo "✓ Provider restarted successfully."

echo -e "\nConsumer 2 retrieving content after Provider restart..."
./bin/consumer -fetch "$CONTENT_ID" -key "$KEY" -out test_recov2.dat -bootstrap "$BOOT_ADDR" -store ./store_client2 > client2.log 2>&1

RECOV2_HASH=$(shasum -a 256 test_recov2.dat | awk '{print $1}')
echo "Consumer 2 Downloaded SHA-256: $RECOV2_HASH"

if [ "$ORIG_HASH" != "$RECOV2_HASH" ]; then
    echo "ERROR: Hash mismatch on Consumer 2 recovery!"
    cat client2.log
    exit 1
fi
echo "✓ SUCCESS: Content persisted and retrievable across provider restart!"

echo -e "\n======================================================================"
echo "    🎉 ALL LIFECYCLE TESTS PASSED: TRUE PROVIDER INDEPENDENCE        "
echo "======================================================================"
