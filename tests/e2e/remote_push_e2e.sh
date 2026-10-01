#!/bin/bash
set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
cd "$ROOT"

echo "======================================================================"
echo "      CIPHER Phase 4: Remote Ingestion & Multi-Provider Replication   "
echo "======================================================================"

rm -rf store_p1 store_p2 store_p3 store_pub store_client1 store_client2
rm -f test_orig.dat test_recovered.dat test_recovered_fault.dat
rm -f bootstrap.log provider1.log provider2.log provider3.log publisher.log client1.log client2.log

export CGO_ENABLED=0

cleanup() {
    echo -e "\nCleaning up background processes..."
    kill $BOOT_PID $PROV1_PID $PROV2_PID $PROV3_PID 2>/dev/null || true
}
trap cleanup EXIT

echo -e "\n[Step 1/6] Building all binaries..."
go build -o bin/bootstrap ./network/cmd/bootstrap
go build -o bin/provider ./nodes/provider
go build -o bin/publisher ./nodes/publisher
go build -o bin/consumer ./nodes/consumer
cp bin/consumer bin/client

echo -e "\n[Step 2/6] Generating 2 MB test payload..."
head -c 2097152 </dev/urandom > test_orig.dat
ORIG_HASH=$(shasum -a 256 test_orig.dat | awk '{print $1}')
echo "Original Payload SHA-256: $ORIG_HASH"

echo -e "\n[Step 3/6] Starting DHT Bootstrap Node..."
./bin/bootstrap -p 48001 -ws-port 0 -identity ./store_pub/boot.key > bootstrap.log 2>&1 &
BOOT_PID=$!
sleep 1

BOOT_ADDR=$(grep "127.0.0.1/tcp/48001/p2p/" bootstrap.log | head -n 1 | awk '{print $NF}')
echo "Bootstrap Multiaddr: $BOOT_ADDR"

echo -e "\n[Step 4/6] Starting 3 Isolated Providers on separate ports & stores..."
./bin/provider -p 48010 -ws-port 0 -identity ./store_p1/p1.key -store ./store_p1 -bootstrap "$BOOT_ADDR" > provider1.log 2>&1 &
PROV1_PID=$!

./bin/provider -p 48020 -ws-port 0 -identity ./store_p2/p2.key -store ./store_p2 -bootstrap "$BOOT_ADDR" > provider2.log 2>&1 &
PROV2_PID=$!

./bin/provider -p 48030 -ws-port 0 -identity ./store_p3/p3.key -store ./store_p3 -bootstrap "$BOOT_ADDR" > provider3.log 2>&1 &
PROV3_PID=$!

sleep 2

P1_ADDR=$(grep "127.0.0.1/tcp/48010/p2p/" provider1.log | head -n 1 | awk '{print $NF}')
P2_ADDR=$(grep "127.0.0.1/tcp/48020/p2p/" provider2.log | head -n 1 | awk '{print $NF}')
P3_ADDR=$(grep "127.0.0.1/tcp/48030/p2p/" provider3.log | head -n 1 | awk '{print $NF}')

echo "Provider 1: $P1_ADDR"
echo "Provider 2: $P2_ADDR"
echo "Provider 3: $P3_ADDR"

echo -e "\n[Step 5/6] Publisher pushes 2 MB file across Providers with Replication R=2..."
./bin/publisher -file test_orig.dat -bootstrap "$BOOT_ADDR" -providers "$P1_ADDR,$P2_ADDR,$P3_ADDR" -replication 2 -push > publisher.log 2>&1

CONTENT_ID=$(grep "ContentID     :" publisher.log | awk '{print $NF}')
KEY=$(grep "Decryption Key:" publisher.log | awk '{print $NF}')

if [ -z "$CONTENT_ID" ] || [ -z "$KEY" ]; then
    echo "ERROR: Publisher failed to push content! Logs:"
    cat publisher.log
    exit 1
fi

echo "Publisher Push Completed Successfully:"
echo "  - ContentID:      $CONTENT_ID"
echo "  - Decryption Key: $KEY"

# Verify publisher has cleanly exited
if pgrep -f "bin/publisher.*test_orig.dat" > /dev/null; then
    echo "ERROR: Publisher process still running after -push!"
    exit 1
fi
echo "✓ Verified: Publisher process exited after satisfying replication invariant."

echo -e "\n[Step 6/6] Consumer 1 discovers providers via DHT and swarm-retrieves content..."
./bin/consumer -fetch "$CONTENT_ID" -key "$KEY" -out test_recovered.dat -bootstrap "$BOOT_ADDR" -store ./store_client1 > client1.log 2>&1

RECOVERED_HASH=$(shasum -a 256 test_recovered.dat | awk '{print $1}')
echo "Consumer 1 Downloaded SHA-256: $RECOVERED_HASH"

if [ "$ORIG_HASH" != "$RECOVERED_HASH" ]; then
    echo "ERROR: Hash mismatch on Consumer 1 recovery!"
    cat client1.log
    exit 1
fi
echo "✓ SUCCESS: Consumer 1 retrieved file via DHT discovery from remote replicated providers!"

echo -e "\n[Fault Tolerance Test] Terminating Provider 1..."
kill $PROV1_PID 2>/dev/null || true
echo "✓ Provider 1 terminated."

echo "Consumer 2 downloading content from remaining Providers (2 & 3)..."
./bin/consumer -fetch "$CONTENT_ID" -key "$KEY" -out test_recovered_fault.dat -bootstrap "$BOOT_ADDR" -store ./store_client2 > client2.log 2>&1

FAULT_RECOVERED_HASH=$(shasum -a 256 test_recovered_fault.dat | awk '{print $1}')
echo "Consumer 2 Downloaded SHA-256: $FAULT_RECOVERED_HASH"

if [ "$ORIG_HASH" != "$FAULT_RECOVERED_HASH" ]; then
    echo "ERROR: Hash mismatch on fault-tolerant recovery!"
    cat client2.log
    exit 1
fi
echo "✓ SUCCESS: Consumer 2 successfully reconstructed file despite dead Provider 1 (R=2 redundancy verified)!"

echo -e "\n======================================================================"
echo "    🎉 ALL PHASE 4 TESTS PASSED: REMOTE INGESTION + REPLICATION      "
echo "======================================================================"
