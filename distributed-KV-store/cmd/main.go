package main

import (
	"distributed-key-value-store/pkg/config"
	"distributed-key-value-store/pkg/raft"
	"distributed-key-value-store/pkg/storage"
	"flag"
	"log"
	"strings"
)

func main() {
	bindAddr := flag.String("bind", ":8080", "Bind address")
	nodeID := flag.String("id", "", "Node ID")
	peers := flag.String("peers", "", "Comma-separated list of peer addresses")
	raftDir := flag.String("raft-dir", "./raft", "Raft data directory")
	dataDir := flag.String("data-dir", "./data", "Data directory")
	flag.Parse()

	if *nodeID == "" {
		log.Fatalf("Node ID is required")
	}

	conf := &config.Config{
		BindAddr: *bindAddr,
		NodeID:   *nodeID,
		Peers:    splitString(*peers),
		RaftDir:  *raftDir,
		DataDir:  *dataDir,
	}

	storage := storage.NewStorage()
	raftNode, err := raft.NewRaftStorage(conf, storage)
	if err != nil {
		log.Fatalf("Failed to create raft node: %v", err)
	}

	if len(conf.Peers) > 0 {
		for _, peer := range conf.Peers {
			if err := raftNode.Join(peer, peer); err != nil {
				log.Fatalf("Failed to join peer %s: %v", peer, err)
			}
		}
	}

	transport.StartServer(conf, storage, raftNode.raft)
}

func splitString(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
