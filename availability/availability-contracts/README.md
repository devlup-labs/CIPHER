# Availability

## Owns

- Availability agreements, commitments, and file availability state.
- Epochs, randomness, challenge generation, and chunk selection.
- Challenge responses, proof verification, freshness, deadlines, and replay protection.
- Availability results, failure handling, and availability-specific smart contracts.

## Does not own

- Peer discovery, Kademlia/DHT implementation, or chunk transport.
- Payment calculation, payment state, escrow, or settlement.
- General reputation systems.

The availability team owns this domain's internal structure, including [`contracts/`](contracts/).
