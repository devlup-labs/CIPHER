# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial production-grade repository structure setup (docs, scripts, Makefile, etc.).
- Provider startup now logs the libp2p peer ID and full provider multiaddrs for easier client connection setup.
- Added a `-quic` CLI flag for both provider and client to enable QUIC transport explicitly.
- Added a shared 30-second P2P operation timeout for host startup, peer connection, chunk requests, mDNS peer connects, and stream handshakes.
- Initial Ethereum integration for provider identity management.
- Added Ethereum wallet generation and persistent key storage for providers.
- Added in-memory provider registry and stake management for the MVP.
- Added challenge creation,resolution framework and stake slashing for provider accountability are MVP implementations.
-  Challenge state is currently maintained in memory and is intended to be replaced by persistent or on-chain challenge management in a future.
- Added Ethereum package unit tests.

### Changed
- QUIC transport is now disabled by default; TCP/WebSocket remain enabled for stable local runs and tests.
- P2P tests now use bounded contexts so failed networking operations do not hang indefinitely.
- - Provider startup now registers its identity with the in-memory Ethereum registry.
- P2P client creates a challenge when a provider fails to reveal the decryption key within the timeout.
- - Providers now attempt to resolve active challenges after successfully revealing the decryption key.

### Fixed
- Fixed local P2P verification by making provider connection details visible in provider logs.
- Fixed the default P2P test path by avoiding the unstable QUIC dependency path unless `-quic` is requested.
