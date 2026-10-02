#!/usr/bin/env bash
set -e

echo "=================================================="
echo "🚀 STARTING LOCAL ANVIL EVM TESTBED"
echo "=================================================="

# 1. Kill any stale anvil processes on port 8545
if lsof -ti :8545 >/dev/null 2>&1; then
    echo "[*] Killing existing process on port 8545..."
    kill -9 $(lsof -ti :8545) 2>/dev/null || true
    sleep 1
fi

# 2. Start Anvil in background
echo "[*] Spawning Anvil node on 127.0.0.1:8545..."
anvil --port 8545 --silent &
ANVIL_PID=$!

# Ensure Anvil is killed when this script exits
cleanup() {
    echo ""
    echo "[*] Stopping Anvil node (PID: $ANVIL_PID)..."
    kill $ANVIL_PID 2>/dev/null || true
}
trap cleanup EXIT

# 3. Wait for Anvil to be ready
echo "[*] Waiting for Anvil RPC endpoint to respond..."
while ! curl -s -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://127.0.0.1:8545 >/dev/null 2>&1; do
    sleep 0.2
done
echo "[+] Anvil is healthy and listening on http://127.0.0.1:8545!"

# 4. Execute Step 1 (Deploy contracts, stake provider, open channel, commit round)
echo ""
echo "=================================================="
echo "📦 STEP 1: DEPLOY & CONFIGURE (BROADCASTING TO ANVIL)"
echo "=================================================="
forge script script/Step1_Setup.s.sol --rpc-url http://127.0.0.1:8545 --broadcast

# 5. Advance block height past confirmation delay
echo ""
echo "=================================================="
echo "⏳ STEP 2: ADVANCING BLOCK HEIGHT (ANVIL MINE)"
echo "=================================================="
echo "[*] Current block: $(cast block-number --rpc-url http://127.0.0.1:8545)"
echo "[*] Mining 20 blocks to satisfy CONFIRMATION_DELAY..."
cast rpc anvil_mine 20 --rpc-url http://127.0.0.1:8545 >/dev/null
echo "[+] New block: $(cast block-number --rpc-url http://127.0.0.1:8545)"

# 6. Execute Step 2 (Settle winning ticket)
echo ""
echo "=================================================="
echo "💰 STEP 3: EXECUTE SETTLEMENT ON ANVIL"
echo "=================================================="
forge script script/Step2_Settle.s.sol --rpc-url http://127.0.0.1:8545 --broadcast

echo ""
echo "=================================================="
echo "✅ ANVIL TEST PASSED SUCCESSFULLY!"
echo "=================================================="
