# kafka-zero

A minimal, from-scratch distributed commit log and message broker inspired by Apache Kafka, written in Go.

---

## Current Status

The core custom binary wire protocol and serialization layer have been designed, implemented, and unit tested.

### Implemented So Far

- **Binary Wire Protocol & Framing (`internal/protocol`)**:
  - Defined request and response message types (`Produce`, `Fetch`, `Metadata`, `CreateTopic`, `Replicate`, `TopicNotification`, `ErrorResponse`).
  - Request encoders with byte-level validation:
    - `EncodeProduceRequest`
    - `EncodeFetchRequest`
    - `EncodeMetadataRequest`
    - `EncodeCreateTopicRequest`
    - `EncodeReplicateRequest`
    - `EncodeTopicNotification`
    - `EncodeErrorResponse`
  - Response decoders with strict bounds and error checking:
    - `DecodeProduceResponse`
    - `DecodeFetchResponse`
    - `DecodeMetadataResponse`
  - Unit test suite covering payload encoding, response parsing, boundary conditions, and corrupt packet rejection.

- **Project Scaffolding**:
  - Base architecture structured into `cmd/broker`, `cmd/client`, and `internal/` packages (`protocol`, `broker`, `client`).
  - Foundational metadata models (`BrokerInfo`, `BrokerMetadata`, `TopicMetadata`, `PartitionMetadata`).

---

## Roadmap (What's Ahead)

1. **Storage Engine (Append-Only Commit Log)**
   - Disk-backed segment files with memory-mapped or sequential file I/O.
   - Sparse offset indexing for fast record lookups.
   - Log retention and segment rolling policies.

2. **Broker Network Layer (`cmd/broker`, `internal/broker`)**
   - TCP server event loop handling framed binary requests concurrently.
   - Request router dispatching to produce, fetch, and metadata handlers.

3. **Topic & Partition Management**
   - In-memory partition state and log directory isolation.
   - Dynamic topic creation and partition routing.

4. **Replication & Cluster Coordination**
   - Leader-follower replication over TCP using `Replicate` RPCs.
   - High watermark tracking and basic leader election/failover.

5. **Client Implementation (`cmd/client`, `internal/client`)**
   - Go client library providing intuitive Producer and Consumer interfaces.
   - CLI tools for publishing and consuming messages from topics.
