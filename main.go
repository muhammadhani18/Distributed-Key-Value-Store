package main

import (
	"context"
	"log"
	"net"

	pb "distributed-kv-store/kvstore"
	"distributed-kv-store/hash"
	"distributed-kv-store/store"

	"google.golang.org/grpc"
)

// Server implements the KeyValueService and includes the hash ring.
type Server struct {
	pb.UnimplementedKeyValueServiceServer
	store       *store.KeyValueStore
	hashRing    *hash.HashRing
	currentNode string
	nodes       []string
}

// Put inserts or updates a key-value pair.
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	// Determine the responsible node for the key.
	targetNode := s.hashRing.GetNode(req.Key)
	if targetNode != s.currentNode {
		// Forward the request to the responsible node via gRPC.
		conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pb.NewKeyValueServiceClient(conn)
		return client.Put(ctx, req)
	}

	// Handle the request locally.
	s.store.Put(req.Key, req.Value)
	return &pb.PutResponse{Success: true}, nil
}

// Get retrieves a value by key.
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	// Determine the responsible node for the key.
	targetNode := s.hashRing.GetNode(req.Key)
	if targetNode != s.currentNode {
		// Forward the request to the responsible node via gRPC.
		conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pb.NewKeyValueServiceClient(conn)
		return client.Get(ctx, req)
	}

	// Handle the request locally.
	value, found := s.store.Get(req.Key)
	return &pb.GetResponse{Value: value, Found: found}, nil
}

// Delete removes a key-value pair.
func (s *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	// Determine the responsible node for the key.
	targetNode := s.hashRing.GetNode(req.Key)
	if targetNode != s.currentNode {
		// Forward the request to the responsible node via gRPC.
		conn, err := grpc.Dial(targetNode, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pb.NewKeyValueServiceClient(conn)
		return client.Delete(ctx, req)
	}

	// Handle the request locally.
	s.store.Delete(req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func main() {
	// Define the nodes in the cluster.
	nodes := []string{"localhost:50051", "localhost:50052", "localhost:50053"} // Example: 3 nodes
	currentNode := "localhost:50051"                                         // Current node address

	// Initialize the hash ring and add all nodes.
	hashRing := hash.NewHashRing(3) // 3 replicas
	for _, node := range nodes {
		hashRing.AddNode(node)
	}

	// Initialize the store and server.
	store := store.NewKeyValueStore()
	server := &Server{
		store:       store,
		hashRing:    hashRing,
		currentNode: currentNode,
		nodes:       nodes,
	}

	// Start the gRPC server.
	grpcServer := grpc.NewServer()
	pb.RegisterKeyValueServiceServer(grpcServer, server)

	listener, err := net.Listen("tcp", currentNode)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", currentNode, err)
	}

	log.Printf("Node %s is listening...", currentNode)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
