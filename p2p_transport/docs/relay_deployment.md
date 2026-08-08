# Relay Node Deployment Guide

The `relay` node provides NAT traversal capabilities (hole punching & circuit relaying) for standard peers in the CIPHER network. 

This guide describes how to deploy the relay node on an Ubuntu server (Linux/amd64).

## 1. Build the Binary
On your development machine or the server, build the binary for Linux AMD64:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/relay-linux-amd64 ./cmd/relay
```

Upload the `bin/relay-linux-amd64` binary to your Ubuntu server.

## 2. Server Configuration
Ensure your server's firewall has the following ports open for inbound traffic:
- **TCP 4001**
- **UDP 4001** (for QUIC)

## 3. Running the Relay

### Option A: Using `tmux` (For quick testing/sessions)
If you want to run the relay in a detached terminal session:
```bash
# Start a new tmux session
tmux new -s cipher-relay

# Run the relay
./relay-linux-amd64

# Detach from the session by pressing Ctrl+B, then D.
```

### Option B: Using `systemd` (Recommended for Production)
To run the relay as a persistent background service, create a systemd service file:

1. Move the binary to a system location:
```bash
sudo mv relay-linux-amd64 /usr/local/bin/cipher-relay
sudo chmod +x /usr/local/bin/cipher-relay
```

2. Create a service file:
```bash
sudo nano /etc/systemd/system/cipher-relay.service
```

3. Paste the following configuration:
```ini
[Unit]
Description=CIPHER P2P Relay Node
After=network.target

[Service]
ExecStart=/usr/local/bin/cipher-relay
Restart=always
User=root
# We run as root or a dedicated user to ensure ~/.config/cipher/ is accessible
Environment="CIPHER_CONFIG_DIR=/etc/cipher"

[Install]
WantedBy=multi-user.target
```

4. Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable cipher-relay
sudo systemctl start cipher-relay
```

## 4. Collecting Connection Information
Once the relay starts, it will output its Peer ID and local Multiaddresses. To allow peers to connect to your relay, you must replace the local IP in the multiaddress with your server's **Public IP**.

**Example output:**
```
Relay Service Started!
Relay Peer ID: 12D3KooWDn...

Relay Multiaddresses (for other peers to connect):
/ip4/127.0.0.1/tcp/4001/p2p/12D3KooWDn...
```

**Your Public Multiaddress:**
If your server's public IP is `203.0.113.50`, your public relay multiaddress is:
`/ip4/203.0.113.50/tcp/4001/p2p/12D3KooWDn...`
