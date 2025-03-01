package raft

import (
	"distributed-key-value-store/pkg/config"
	"distributed-key-value-store/pkg/storage"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type RaftNode struct {
	raft *raft.Raft
}

func NewRaftStorage(config *config.Config, storage *storage.Storage) (*RaftNode, error) {
	logStore, err := raftboltdb.NewBoltStore(config.RaftDir + "/raft.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create log store: %w", err)
	}

	boltDB, err := raftboltdb.NewBoltStore(config.RaftDir + "/raft.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create boltdb: %w", err)
	}

	transport, err := raft.NewTCPTransport(config.BindAddr, nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create TCP transport: %w", err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(config.RaftDir, 1, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot store: %w", err)
	}

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(config.NodeID)

	raftNode, err := raft.NewRaft(raftConfig, NewFSM(storage), logStore, boltDB, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create raft node: %w", err)
	}

	return &RaftNode{raft: raftNode}, nil
}

func (r *RaftNode) Join(nodeID, addr string) error {
	configFuture := r.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return fmt.Errorf("failed to get raft configuration: %w", err)
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) {
			// Node already exists, return without adding again
			if srv.Address == raft.ServerAddress(addr) {
				return nil
			}
			// Remove existing node if the address is different
			err := r.raft.RemoveServer(srv.ID, 0, 0).Error()
			if err != nil {
				return fmt.Errorf("failed to remove existing node: %w", err)
			}
		}
	}

	// Adding the new node as a voter
	if err := r.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0).Error(); err != nil {
		return fmt.Errorf("failed to add voter: %w", err)
	}

	return nil
}

func (r *RaftNode) GetLeader() raft.ServerAddress {
	return r.raft.Leader()
}
