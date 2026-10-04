#!/usr/bin/env bash
# ==============================================================================
# CIPHER QUICKSTART: ZERO-CONFIG MULTI-LAPTOP LIVE RUNNER (AZURE CONNECTED)
# ==============================================================================
# Usage:
#   ./quickstart.sh provider [port] [store_dir]
#   ./quickstart.sh publisher <file_path>
#   ./quickstart.sh consumer <content_id> <key> [output_path]
# ==============================================================================

set -e

# Disable Git Bash MSYS automatic POSIX-to-Windows path conversion for multiaddrs
export MSYS_NO_PATHCONV=1
export MSYS2_ARG_CONV_EXCL="*"

# Windows vs POSIX executable suffix
EXE=""
if [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "cygwin"* || "$OSTYPE" == "win32"* ]]; then
    EXE=".exe"
fi

# Azure Infrastructure Coordinates
AZURE_IP="${AZURE_IP:-20.197.30.171}"
BOOTSTRAP_PEER_ID="${AZURE_BOOTSTRAP_PEER_ID:-12D3KooWBzG7KT4Z7sNzyzwWGBowEpiwvPMcfgRnf8XWUtkC6PcJ}"
RELAY_PEER_ID="${AZURE_RELAY_PEER_ID:-12D3KooWSwofcWofnKhvEQNfE8YD4Xv4HosdXGRpJ334xkyiwvbT}"

BOOTSTRAP="/ip4/${AZURE_IP}/tcp/4003/p2p/${BOOTSTRAP_PEER_ID}"
RELAY="/ip4/${AZURE_IP}/tcp/4001/p2p/${RELAY_PEER_ID}"
ETH_RPC="http://${AZURE_IP}:8545"
ENTROPY="0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9"
PROVIDER_ETH="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CLIENT_ETH_KEY="5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"

ROLE="${1:-}"

mkdir -p bin

case "$ROLE" in
  provider|prov)
    PORT="${2:-4101}"
    STORE="${3:-./store_p_${PORT}}"
    echo "======================================================================"
    echo " 🚀 STARTING CIPHER STORAGE PROVIDER"
    echo "    Port:      $PORT"
    echo "    Store:     $STORE"
    echo "    Bootstrap: $BOOTSTRAP"
    echo "======================================================================"
    go build -o "bin/provider${EXE}" ./nodes/provider
    chmod +x "bin/provider${EXE}" 2>/dev/null || true
    exec "./bin/provider${EXE}" -p "$PORT" -store "$STORE" \
      -bootstrap "$BOOTSTRAP" \
      -relay "$RELAY" \
      --eth-rpc "$ETH_RPC" \
      --entropy-addr "$ENTROPY"
    ;;

  publisher|pub)
    FILE="${2:-}"
    if [ -z "$FILE" ]; then
      echo "❌ Error: Please provide a file to publish."
      echo "Usage: ./quickstart.sh publisher <path_to_file>"
      echo "Example: head -c 1048576 </dev/urandom > test.dat && ./quickstart.sh publisher test.dat"
      exit 1
    fi
    if [ ! -f "$FILE" ]; then
      echo "❌ Error: File '$FILE' does not exist."
      exit 1
    fi
    echo "======================================================================"
    echo " 📤 INGESTING & PUBLISHING CONTENT (Replication R=2)"
    echo "    File:      $FILE"
    echo "    Bootstrap: $BOOTSTRAP"
    echo "======================================================================"
    go build -o "bin/publisher${EXE}" ./nodes/publisher
    chmod +x "bin/publisher${EXE}" 2>/dev/null || true
    "./bin/publisher${EXE}" -file "$FILE" -replication 2 -push \
      -bootstrap "$BOOTSTRAP" \
      -relay "$RELAY"
    ;;

  consumer|client)
    CID="${2:-}"
    KEY="${3:-}"
    OUT="${4:-./downloaded.dat}"
    if [ -z "$CID" ] || [ -z "$KEY" ]; then
      echo "❌ Error: Missing ContentID or Decryption Key."
      echo "Usage: ./quickstart.sh consumer <ContentID> <DecryptionKey> [output_file]"
      exit 1
    fi
    echo "======================================================================"
    echo " 📥 DOWNLOADING CONTENT (DHT Discovery & EIP-712 Micro-payments)"
    echo "    ContentID: $CID"
    echo "    Output:    $OUT"
    echo "    Bootstrap: $BOOTSTRAP"
    echo "======================================================================"
    go build -o "bin/consumer${EXE}" ./nodes/consumer
    chmod +x "bin/consumer${EXE}" 2>/dev/null || true
    "./bin/consumer${EXE}" -fetch "$CID" -key "$KEY" -out "$OUT" \
      -bootstrap "$BOOTSTRAP" \
      -relay "$RELAY" \
      --eth-rpc "$ETH_RPC" \
      --eth-key "$CLIENT_ETH_KEY" \
      --entropy-addr "$ENTROPY" \
      --provider-eth-addr "$PROVIDER_ETH"
    echo ""
    echo "[✓] Download complete: $OUT"
    if command -v shasum >/dev/null 2>&1; then
      echo "SHA-256: $(shasum -a 256 "$OUT" | awk '{print $1}')"
    elif command -v sha256sum >/dev/null 2>&1; then
      echo "SHA-256: $(sha256sum "$OUT" | awk '{print $1}')"
    fi
    ;;

  *)
    echo "======================================================================"
    echo "               CIPHER ZERO-CONFIG RUNNER (AZURE VM)"
    echo "======================================================================"
    echo "Usage:"
    echo "  1. Run Storage Provider:  ./quickstart.sh provider [port] [store_dir]"
    echo "  2. Publish File:          ./quickstart.sh publisher <file_path>"
    echo "  3. Download as Consumer:  ./quickstart.sh consumer <content_id> <key> [output_path]"
    echo "======================================================================"
    exit 1
    ;;
esac
