# Distributed Key-Value Store

A fully functional **Distributed Key-Value Store** implemented in **Go** using **gRPC**. It supports basic CRUD operations (`PUT`, `GET`, and `DELETE`) with consistent hashing to distribute keys across multiple nodes. This project demonstrates distributed systems concepts, focusing on horizontal scalability, consistency, and efficient data distribution.

---

## Features
- **Distributed System**: Keys are evenly distributed across nodes using **consistent hashing**.
- **gRPC for Communication**: Nodes communicate with each other using gRPC.
- **Node-to-Node Forwarding**: Requests are forwarded to the responsible node if necessary.
- **Thread-Safe In-Memory Storage**: Local storage uses a thread-safe map for concurrency.
- **Horizontal Scalability**: Easily add nodes to the cluster by modifying the hash ring.

---

## Project Structure

distributed-kv-store/
├── kvstore/
│     ├── kvstore.proto          # Protocol Buffers definition
│     ├── kvstore.pb.go          # Generated Protobuf code
│     ├── kvstore_grpc.pb.go     # Generated gRPC code
├── store/
│     ├── kvstore.go             # In-memory Key-Value Store logic
├── hash/
│     ├── hash_ring.go           # Consistent Hashing logic
├── main.go                      # gRPC server with node-to-node communication
├── go.mod                       # Go module file
└── README.md                    # Documentation


---

## How it Works

1. **Consistent Hashing**: 
   - Distributes keys across nodes using consistent hashing, ensuring even distribution and minimizing rehashing when nodes are added or removed.

2. **gRPC for Communication**:
   - Handles communication between clients and nodes, as well as between nodes for forwarding requests.

3. **Node-to-Node Communication**:
   - Requests for keys not owned by the current node are forwarded to the correct node.

4. **Thread-Safe Local Storage**:
   - Each node uses a thread-safe in-memory `map` with read/write locks to store key-value pairs.

---

## Setup and Installation

### Prerequisites
1. Install **Go** (1.20+): [https://go.dev/dl/](https://go.dev/dl/)
2. Install **Protocol Buffers (protoc)**: [https://protobuf.dev/downloads/](https://protobuf.dev/downloads/)
3. Install the Go gRPC plugins:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

### Clone the repository
```bash
git clone https://github.com/muhammadhani18/distributed-kv-store.git
```

### Generate gRPC Code
Make sure protoc is in your PATH. Then run:
```bash
protoc --go_out=. --go-grpc_out=. kvstore.proto
```

This will generate the following files:

1. kvstore.pb.go: Protobuf message definitions.
2. kvstore_grpc.pb.go: gRPC service stubs.

### Install Dependencies

```bash
go mod tidy
```

## Run the Server

```bash
go run main.go
```

### Node 2
Edit the ```currentNode``` value in ```main.go``` to ```localhost:50052```:

```go
currentNode := "localhost:50052"
```

###Step 2: Test the System
You can test the distributed key-value store using BloomRPC, grpcurl, or a custom client.

#### Using BloomRPC
1. Import the kvstore.proto file into BloomRPC.
2. Set the server address (e.g., localhost:50051).
3. Use the Put, Get, and Delete methods to interact with the key-value store.
4. Using grpcurl
5. Install grpcurl:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Run the following commands to interact with the server:

Put a key-value pair:

```bash
grpcurl -plaintext -d '{"key": "mykey", "value": "myvalue"}' localhost:50051 pb.KeyValueService.Put
```
Get the value of a key:

```bash
grpcurl -plaintext -d '{"key": "mykey"}' localhost:50051 pb.KeyValueService.Get
```

Delete a key:

```bash
grpcurl -plaintext -d '{"key": "mykey"}' localhost:50051 pb.KeyValueService.Delete
```

---

## Code Walkthrough

### Local In-Memory Store (store/kvstore.go)
A thread-safe map is used to store key-value pairs:
```go
type KeyValueStore struct {
    mu    sync.RWMutex
    store map[string]string
}

// Put stores a key-value pair.
func (kv *KeyValueStore) Put(key, value string) {
    kv.mu.Lock()
    defer kv.mu.Unlock()
    kv.store[key] = value
}

// Get retrieves a value by key.
func (kv *KeyValueStore) Get(key string) (string, bool) {
    kv.mu.RLock()
    defer kv.mu.RUnlock()
    value, found := kv.store[key]
    return value, found
}

// Delete removes a key-value pair.
func (kv *KeyValueStore) Delete(key string) {
    kv.mu.Lock()
    defer kv.mu.Unlock()
    delete(kv.store, key)
}

```
### Consistent Hashing (hash/hash_ring.go)
Used to evenly distribute keys across nodes.
```go
type HashRing struct {
    nodes    []string
    replicas int
    ring     map[int]string
    sorted   []int
}

// AddNode adds a new node to the ring.
func (h *HashRing) AddNode(node string) { ... }

// GetNode returns the node responsible for a given key.
func (h *HashRing) GetNode(key string) string { ... }

```

### gRPC Server (main.go)
Handles gRPC requests and forwards them to the appropriate node.

```go
type Server struct {
    store       *store.KeyValueStore
    hashRing    *hash.HashRing
    currentNode string
    nodes       []string
}

func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) { ... }
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) { ... }
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) { ... }

```
---

## Future Improvements
- **Persistent Storage**: Implement a persistent storage solution (e.g., Redis) for data durability.
- **Leader Election**: Add a leader election mechanism to handle node failures.
- **Monitoring and Metrics**: Implement monitoring and metrics collection for the distributed system.
- **Load Balancing**: Implement load balancing to distribute requests evenly across nodes.
- **Encryption**: Add encryption for data in transit and at rest.
- **Authentication**: Implement authentication and authorization for the gRPC server.

