package raft

import (
	"bytes"
	"distributed-key-value-store/pkg/storage"
	"encoding/gob"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/raft"
)

type FSM struct {
	storage *storage.Storage
}

func NewFSM(storage *storage.Storage) *FSM {
	return &FSM{
		storage: storage,
	}
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	var command string
	if err := gob.NewDecoder(bytes.NewReader(log.Data)).Decode(&command); err != nil {
		panic(err)
	}

	parts := strings.Split(command, " ")
	if len(parts) < 2 {
		panic(fmt.Sprintf("invalid command: %s", command))
	}

	switch parts[0] {
	case "SET":
		if len(parts) < 3 {
			panic("invalid SET command")
		}
		key, value := parts[1], parts[2]
		f.storage.Set(key, value)
	case "DELETE":
		key := parts[1]
		f.storage.Delete(key)
	default:
		panic(fmt.Sprintf("unknown command: %s", command))
	}

	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{storage: f.storage}, nil
}

func (f *FSM) Restore(snapshot io.ReadCloser) error {
	defer snapshot.Close()
	dec := gob.NewDecoder(snapshot)
	var db map[string]string
	if err := dec.Decode(&db); err != nil {
		return err
	}

	f.storage.MU.Lock()
	defer f.storage.MU.Unlock()

	f.storage.DB = make(map[string]string)
	for k, v := range db {
		f.storage.DB[k] = v
	}

	return nil
}

type fsmSnapshot struct {
	storage *storage.Storage
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	f.storage.MU.RLock()
	defer f.storage.MU.RUnlock()

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(f.storage.DB); err != nil {
		sink.Cancel()
		return err
	}

	if _, err := sink.Write(buf.Bytes()); err != nil {
		sink.Cancel()
		return err
	}

	return sink.Close()
}

func (f *fsmSnapshot) Release() {}
