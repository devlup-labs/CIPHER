#!/bin/bash
set -e

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"
cd "$ROOT"

echo "======================================================================"
echo "      CIPHER E2E: P2P Chunk Transfer with Live Anvil Payment Settlement"
echo "======================================================================"

rm -rf store_p1 store_pub store_client test_orig.dat test_recovered.dat
rm -f provider.log publisher.log consumer.log anvil.log
mkdir -p bin

# 1. Kill any existing Anvil on port 8545
if lsof -ti :8545 >/dev/null 2>&1; then
    echo "[*] Killing existing process on port 8545..."
    kill -9 $(lsof -ti :8545) 2>/dev/null || true
    sleep 1
fi

# 2. Start Anvil in background
echo "[*] Launching Anvil EVM node on 127.0.0.1:8545..."
anvil --port 8545 --silent > anvil.log 2>&1 &
ANVIL_PID=$!

cleanup() {
    echo -e "\n[*] Cleaning up background processes..."
    kill $PROV_PID $ANVIL_PID 2>/dev/null || true
}
trap cleanup EXIT

# 3. Wait for Anvil to be healthy
while ! curl -s -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://127.0.0.1:8545 >/dev/null 2>&1; do
    sleep 0.2
done
echo "[+] Anvil is listening and healthy!"

# 4. Deploy contracts & configure channels on Anvil
echo -e "\n[Step 1/6] Deploying Smart Contracts & Staking Provider on Anvil..."
(cd payments && forge script script/Step1_Setup.s.sol --rpc-url http://127.0.0.1:8545 --broadcast > /dev/null)

ENTROPY_ADDR="0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
PROVIDER_ETH_ADDR="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CLIENT_ETH_KEY="5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"
CLIENT_ETH_ADDR="0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"

echo "  - EntropySource:    $ENTROPY_ADDR"
echo "  - Provider Wallet:  $PROVIDER_ETH_ADDR"
echo "  - Client Wallet:    $CLIENT_ETH_ADDR"

# 5. Build Go binaries
echo -e "\n[Step 2/6] Building Go binaries (provider, publisher, consumer)..."
go build -o bin/provider ./nodes/provider
go build -o bin/publisher ./nodes/publisher
go build -o bin/consumer ./nodes/consumer

# 6. Start Provider with payment verification
echo -e "\n[Step 3/6] Starting Storage Provider with EVM Ticket Verification..."
./bin/provider -p 49010 -ws-port 0 -identity ./store_p1/p1.key -store ./store_p1 \
  --eth-rpc http://127.0.0.1:8545 --entropy-addr "$ENTROPY_ADDR" > provider.log 2>&1 &
PROV_PID=$!
sleep 2

PROV_ADDR=$(grep "127.0.0.1/tcp/49010/p2p/" provider.log | head -n 1 | awk '{print $NF}')
echo "Provider Multiaddr: $PROV_ADDR"

# 7. Generate test payload and push to Provider
echo -e "\n[Step 4/6] Ingesting and Pushing 128 KiB test payload to Provider..."
head -c 131072 </dev/urandom > test_orig.dat
ORIG_HASH=$(shasum -a 256 test_orig.dat | awk '{print $1}')
echo "Original SHA-256: $ORIG_HASH"

./bin/publisher -file test_orig.dat -providers "$PROV_ADDR" -push > publisher.log 2>&1
CONTENT_ID=$(grep "ContentID     :" publisher.log | awk '{print $NF}')
KEY=$(grep "Decryption Key:" publisher.log | awk '{print $NF}')

echo "  - ContentID:      $CONTENT_ID"
echo "  - Decryption Key: $KEY"

# 8. Consumer downloads file while generating and streaming EIP-712 payment tickets
echo -e "\n[Step 5/6] Consumer downloading chunks with EIP-712 Micro-payment Tickets..."
./bin/consumer -p 49020 -ws-port 0 -identity ./store_client/client.key -store ./store_client \
  -d "$PROV_ADDR" -fetch "$CONTENT_ID" -key "$KEY" -out test_recovered.dat \
  --eth-rpc http://127.0.0.1:8545 --eth-key "$CLIENT_ETH_KEY" \
  --entropy-addr "$ENTROPY_ADDR" --provider-eth-addr "$PROVIDER_ETH_ADDR" > consumer.log 2>&1

RECOVERED_HASH=$(shasum -a 256 test_recovered.dat | awk '{print $1}')
echo "Recovered SHA-256: $RECOVERED_HASH"

if [ "$ORIG_HASH" != "$RECOVERED_HASH" ]; then
    echo "❌ FAILED: Hash mismatch between original and recovered data!"
    exit 1
fi
echo "[✓] Data integrity confirmed bit-for-bit!"

# Verify provider received and verified payment tickets
TICKET_COUNT=$(grep -c "Received valid ticket" provider.log || true)
echo "[✓] Provider verified and recorded $TICKET_COUNT ticket(s) over the wire protocol."
if [ "$TICKET_COUNT" -eq 0 ]; then
    echo "❌ FAILED: No valid payment tickets were recorded by the provider!"
    cat provider.log
    exit 1
fi

# 9. Advance blocks & Execute On-chain Settlement on Anvil
echo -e "\n[Step 6/6] Mining 20 blocks on Anvil and Executing On-chain Settlement..."
cast rpc anvil_mine 20 --rpc-url http://127.0.0.1:8545 > /dev/null

(cd payments && forge script script/Step2_Settle.s.sol --rpc-url http://127.0.0.1:8545 --broadcast)

echo -e "\n======================================================================"
echo "🎉 SUCCESS: FULL P2P WIRE TRANSFER & ANVIL ON-CHAIN SETTLEMENT COMPLETE!"
echo "======================================================================"
